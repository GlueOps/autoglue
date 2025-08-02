package cmd

import (
	"context"
	"fmt"
	"github.com/glueops/autoglue/api"
	"github.com/glueops/autoglue/internal/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	bindPort    string
	bindAddress string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as daemon and serve HTTP/HTTPS request",
	Run: func(cmd *cobra.Command, args []string) {
		addr := fmt.Sprintf("%s:%s", viper.GetString("bind_address"), viper.GetString("bind_port"))

		db.Connect()

		server := api.NewServer(addr)

		fmt.Println("starting server at http://" + addr)

		// Channel to listen for interrupt signals
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Run server in goroutine
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		}()

		// Wait for interrupt
		<-stop
		fmt.Println("shutting down server...")

		// Gracefully shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}

		fmt.Println("server shutdown complete")

	},
}

func init() {
	serveCmd.Flags().StringVar(&bindPort, "bind-port", "8080", "HTTP/HTTPS bind port")
	serveCmd.Flags().StringVar(&bindAddress, "bind-address", "127.0.0.1", "HTTP/HTTPS bind address")

	// Bind flags to Viper keys
	_ = viper.BindPFlag("bind_port", serveCmd.Flags().Lookup("bind-port"))
	_ = viper.BindPFlag("bind_address", serveCmd.Flags().Lookup("bind-address"))

	rootCmd.AddCommand(serveCmd)
}
