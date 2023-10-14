// Package command 命令行
package command

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/thoohv5/person/internal/util"
)

var (
	dir     string
	version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "demo-interface",
	Short: "demo进程",
	Long:  `测试使用`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	version = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dir, "conf", util.AbPath("../../configs/"), "config path, eg: --conf=config.yaml")
}
