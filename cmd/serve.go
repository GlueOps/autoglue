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

	"github.com/glueops/autoglue/internal/api"
	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/config"
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

		_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
		go func() {
			t := time.NewTicker(60 * time.Second)
			defer t.Stop()
			for range t.C {
				_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
			}
		}()

		r := api.NewRouter(rt.DB)

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
