package cli

import (
	"github.com/glueops/autoglue/internal/config"
	"github.com/spf13/cobra"
)

var showConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "Show the current configuration",
	Long:  "Show the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.DebugPrintConfig()
	},
}

func init() {
	rootCmd.AddCommand(showConfigCmd)
}
