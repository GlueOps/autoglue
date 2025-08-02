package cmd

import (
	"fmt"

	"github.com/earthboundkid/versioninfo/v2"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Know the installed version of autoglue",
	Long:  `This command will help you to know the installed version of autoglue`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("autoglue version:", versioninfo.Short())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

}
