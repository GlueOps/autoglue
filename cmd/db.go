package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database utilities",
}

var dbPsqlCmd = &cobra.Command{
	Use:   "psql",
	Short: "Open a psql session to the app database",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.DbURL == "" {
			return errors.New("database.url is empty")
		}
		psql := "psql"
		if runtime.GOOS == "windows" {
			psql = "psql.exe"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 72*time.Hour)
		defer cancel()

		psqlCmd := exec.CommandContext(ctx, psql, cfg.DbURL)
		psqlCmd.Stdin, psqlCmd.Stdout, psqlCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		fmt.Println("Launching psql…")
		return psqlCmd.Run()
	},
}

func init() {
	dbCmd.AddCommand(dbPsqlCmd)

	rootCmd.AddCommand(dbCmd)
}
