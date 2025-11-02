package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autoglue",
	Short: "Autoglue Kubernetes Cluster Management",
	Long:  "autoglue is used to manage the lifecycle of kubernetes clusters on GlueOps supported cloud providers",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := serveCmd.RunE(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_ = cmd.Help()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize()
}
