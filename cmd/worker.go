package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the background job worker",
	Long: "Run River workers and the periodic job schedule. This process serves " +
		"no HTTP traffic. Run `autoglue serve` for the API.\n\n" +
		"Multiple replicas are safe: River leases each job to one worker, and " +
		"the periodic schedule runs on the elected leader only.",
	RunE: func(_ *cobra.Command, _ []string) error {
		rt := app.NewRuntime()
		defer rt.Close()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		jobs, err := bg.NewWorkerClient(rt.Pool, bg.Deps{
			DB:      rt.DB,
			BaseURL: rt.Cfg.BaseURL,
		})
		if err != nil {
			return fmt.Errorf("init river client: %w", err)
		}

		if err := jobs.Start(ctx); err != nil {
			return fmt.Errorf("start workers: %w", err)
		}

		fmt.Println("🛠  Worker running (queues: default, maintenance, clusters)")

		<-ctx.Done()
		fmt.Println("\n⏳ Draining jobs...")

		// Stop waits for in-flight jobs to finish. Anything still running when
		// this budget expires is left for another worker to pick up after the
		// rescuer notices the lease has lapsed.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := jobs.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("stop workers: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
