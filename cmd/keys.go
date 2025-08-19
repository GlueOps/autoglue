package cmd

import (
	"github.com/glueops/autoglue/internal/db"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage autoglue encryption keys",
	Long:  "Manage autoglue master encryption keys used for securing data.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if db.DB != nil {
			return nil
		}
		db.Connect()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(rotateMasterCmd)
	keysCmd.AddCommand(createMasterCmd)
}
