package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dyaksa/archer"
	"github.com/glueops/autoglue/internal/api"
	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/web"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	RunE: func(_ *cobra.Command, _ []string) error {
		rt := app.NewRuntime()

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var pgwebInst *web.Pgweb

		jobs, err := bg.NewJobs(rt.DB, cfg.DbURL)
		if err != nil {
			log.Fatalf("failed to init background jobs: %v", err)
		}

		// Start workers in background ONCE
		go func() {
			if err := jobs.Start(); err != nil {
				log.Fatalf("failed to start background jobs: %v", err)
			}
		}()
		defer jobs.Stop()

		// daily cleanups
		{
			// schedule next 03:30 local time
			next := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Add(3*time.Hour + 30*time.Minute)
			_, _ = jobs.Enqueue(
				context.Background(),
				uuid.NewString(),
				"archer_cleanup",
				bg.CleanupArgs{RetainDays: 7, Table: "jobs"},
				archer.WithScheduleTime(next),
				archer.WithMaxRetries(1),
			)

			// schedule next 03:45 local time
			next2 := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Add(3*time.Hour + 45*time.Minute)
			_, _ = jobs.Enqueue(
				context.Background(),
				uuid.NewString(),
				"tokens_cleanup",
				bg.TokensCleanupArgs{},
				archer.WithScheduleTime(next2),
				archer.WithMaxRetries(1),
			)

			_, _ = jobs.Enqueue(
				context.Background(),
				uuid.NewString(),
				"db_backup_s3",
				bg.DbBackupArgs{},
				archer.WithMaxRetries(1),
				archer.WithScheduleTime(time.Now().Add(1*time.Hour)),
			)
		}

		// Periodic scheduler
		schedCtx, schedCancel := context.WithCancel(context.Background())
		defer schedCancel()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ticker.C:
					_, err := jobs.Enqueue(
						context.Background(),
						uuid.NewString(),
						"bootstrap_bastion",
						bg.BastionBootstrapArgs{},
						archer.WithMaxRetries(3),
						// while debugging, avoid extra schedule delay:
						archer.WithScheduleTime(time.Now().Add(60*time.Second)),
					)
					if err != nil {
						log.Printf("failed to enqueue bootstrap_bastion: %v", err)
					}
				case <-schedCtx.Done():
					return
				}
			}
		}()

		_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
		go func() {
			t := time.NewTicker(60 * time.Second)
			defer t.Stop()
			for range t.C {
				_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
			}
		}()

		var studioHandler http.Handler
		r := api.NewRouter(rt.DB, jobs, nil)

		if cfg.DBStudioEnabled {
			dbURL := cfg.DbURLRO
			if dbURL == "" {
				dbURL = cfg.DbURL
			}

			pgwebInst, err = web.StartPgweb(
				dbURL,
				cfg.DBStudioBind,
				cfg.DBStudioPort,
				true,
				cfg.DBStudioUser,
				cfg.DBStudioPass,
			)
			if err != nil {
				log.Printf("pgweb failed to start: %v", err)
			} else {
				studioHandler = http.HandlerFunc(pgwebInst.Proxy())
				r = api.NewRouter(rt.DB, jobs, studioHandler)
				log.Printf("pgweb running on http://%s:%s (proxied at /db-studio/)", cfg.DBStudioBind, pgwebInst.Port())
			}
		}

		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

		srv := &http.Server{
			Addr:         addr,
			Handler:      TimeoutExceptUpgrades(r, 60*time.Second, "request timed out"), // global safety
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			fmt.Printf("🚀 API running on http://%s (ui.dev=%v)\n", addr, cfg.UIDev)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("server error: %v", err)
			}
		}()

		<-ctx.Done()
		fmt.Println("\n⏳ Shutting down...")
		if pgwebInst != nil {
			_ = pgwebInst.Stop(context.Background())
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func TimeoutExceptUpgrades(next http.Handler, d time.Duration, msg string) http.Handler {
	timeout := http.TimeoutHandler(next, d, msg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If this is an upgrade (e.g., websocket), don't wrap.
		if isUpgrade(r) {
			next.ServeHTTP(w, r)
			return
		}
		timeout.ServeHTTP(w, r)
	})
}

func isUpgrade(r *http.Request) bool {
	// Connection: Upgrade, Upgrade: websocket
	if strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
		return true
	}
	return false
}
