package cmd

import (
	"github.com/glueops/autoglue/internal/config"
	"github.com/spf13/cobra"
)

var showConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "Display loaded configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.Load()
		config.DebugPrintConfig()
	},
}

func init() {
	rootCmd.AddCommand(showConfigCmd)
}
