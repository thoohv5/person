package command

import (
	"flag"

	"github.com/spf13/cobra"

	"github.com/thoohv5/person/internal/constant"
)

// serverCmd represents the base command when called without any subcommands
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务器",
	Long:  `HTTP服务器`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()

		// 框架服务
		app, cleanup, err := initApp(constant.ConfigPath(dir))
		if err != nil {
			panic(err)
		}
		defer cleanup()

		// start and wait for stop signal
		if rErr := app.Run(); rErr != nil {
			panic(rErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
