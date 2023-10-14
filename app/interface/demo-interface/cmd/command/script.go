package command

import (
	"context"
	"flag"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/thoohv5/person/internal/constant"
)

// versionCmd represents the base command when called without any subcommands
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "定时任务",
	Long: `This program runs command on the cron. Supported commands are:
  - xxx - cron script name`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()

		// 框架服务
		app, cleanup, err := initCron(constant.ConfigPath(dir))
		if err != nil {
			panic(err)
		}
		defer cleanup()

		if len(args) == 0 {
			fmt.Println("亲输入执行的脚本名称")
			return
		}
		for _, cron := range app {
			if rt := reflect.Indirect(reflect.ValueOf(cron)).Type(); rt.Name() == args[0] {
				cron.Run(context.Background())()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)
}
