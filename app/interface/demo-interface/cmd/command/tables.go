package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thoohv5/person/app/interface/demo-interface/migrations"
)

// versionCmd represents the base command when called without any subcommands
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "表名",
	Long:  `可执行程序的数据库中所有表名`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tablenames: %s\n", migrations.GetTables())
	},
}

func init() {
	rootCmd.AddCommand(tablesCmd)
}
