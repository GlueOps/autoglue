package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autoglue [command]",
	Short: "autoglue is used to manage the lifecycle of kubernetes clusters on GlueOps supported cloud providers",
	Long:  "autoglue is used to manage the lifecycle of kubernetes clusters on GlueOps supported cloud providers",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			serveCmd.Run(cmd, []string{})
		} else {
			_ = cmd.Help()
		}
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func checkNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
