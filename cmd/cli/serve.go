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

	"github.com/glueops/autoglue/internal/api"
	"github.com/glueops/autoglue/internal/db"
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

		// Resolve bind address/port from viper (flags/env/config/defaults)
		addr := fmt.Sprintf("%s:%s", viper.GetString("bind_address"), viper.GetString("bind_port"))

		// Build server (uses Chi router inside)
		srv := api.NewServer(addr)

		// Start server
		errCh := make(chan error, 1)
		go func() {
			log.Printf("HTTP server listening on http://%s (ui.dev=%v)", addr, viper.GetBool("ui.dev"))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
			close(errCh)
		}()

		// Handle OS signals for graceful shutdown
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
	// Flags to override bind address/port
	serveCmd.Flags().StringVar(&bindAddress, "bind-address", "", "Address to bind the HTTP server (default 127.0.0.1)")
	serveCmd.Flags().StringVar(&bindPort, "bind-port", "", "Port to bind the HTTP server (default 8080)")

	// Bind flags to viper keys
	_ = viper.BindPFlag("bind_address", serveCmd.Flags().Lookup("bind-address"))
	_ = viper.BindPFlag("bind_port", serveCmd.Flags().Lookup("bind-port"))

	// Register command
	rootCmd.AddCommand(serveCmd)
}
