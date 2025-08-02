package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autoglue [command]",
	Short: "autoglue is used to manage the lifecycle of kubernetes clusters on GlueOps supported cloud providers",
	Long:  "autoglue is used to manage the lifecycle of kubernetes clusters on GlueOps supported cloud providers",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		checkNilErr(err)
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func checkNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
