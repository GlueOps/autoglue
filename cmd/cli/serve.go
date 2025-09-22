package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dyaksa/archer"
	"github.com/glueops/autoglue/internal/api"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/db"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	bindPort    string
	bindAddress string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		jobs, err := bg.NewJobs()
		if err != nil {
			log.Fatalf("failed to init background jobs: %v", err)
		}

		// Start workers in background ONCE
		go func() {
			if err := jobs.Start(); err != nil {
				log.Fatalf("failed to start background workers: %v", err)
			}
		}()
		defer jobs.Stop()

		// Enqueue one job immediately
		/*
			id := uuid.NewString()
			if _, err := jobs.Enqueue(context.Background(), id, "bootstrap_bastion", bg.BastionBootstrapArgs{}); err != nil {
				log.Fatalf("[enqueue] failed (job_id=%s): %v", id, err)
			}
			log.Printf("[enqueue] queued (job_id=%s)", id)

			// Verify the row exists
			if got, err := jobs.Client.Get(context.Background(), id); err != nil {
				log.Fatalf("[verify] Get failed (job_id=%s): %v", id, err)
			} else if j, ok := got.(*job.Job); ok {
				log.Printf("[verify] Get ok (job_id=%s, status=%s)", j.ID, j.Status)
			} else {
				log.Printf("[verify] Get ok (job_id=%s) but unexpected type %T", id, got)
			}
		*/
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
						archer.WithScheduleTime(time.Now().Add(10*time.Second)),
					)
					if err != nil {
						log.Printf("failed to enqueue bootstrap_bastion: %v", err)
					}
				case <-schedCtx.Done():
					return
				}
			}
		}()

		// HTTP server
		addr := fmt.Sprintf("%s:%s", viper.GetString("bind_address"), viper.GetString("bind_port"))
		srv := api.NewServer(addr)

		errCh := make(chan error, 1)
		go func() {
			log.Printf("HTTP server listening on http://%s (ui.dev=%v)", addr, viper.GetBool("ui.dev"))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
			close(errCh)
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-stop:
			log.Printf("Received signal: %s — shutting down...", sig)
		case err := <-errCh:
			if err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}

		schedCancel()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v; forcing close", err)
			_ = srv.Close()
		} else {
			log.Println("Server stopped cleanly.")
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&bindAddress, "bind-address", "", "Address to bind the HTTP server (default 127.0.0.1)")
	serveCmd.Flags().StringVar(&bindPort, "bind-port", "", "Port to bind the HTTP server (default 8080)")
	_ = viper.BindPFlag("bind_address", serveCmd.Flags().Lookup("bind-address"))
	_ = viper.BindPFlag("bind_port", serveCmd.Flags().Lookup("bind-port"))
	rootCmd.AddCommand(serveCmd)
}
