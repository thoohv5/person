package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本",
	Long:  `可执行程序的版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
