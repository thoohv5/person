package command

import (
	"github.com/spf13/cobra"

	"github.com/thoohv5/person/app/interface/demo-interface/boot"
	im "github.com/thoohv5/person/app/interface/demo-interface/migrations"
	"github.com/thoohv5/person/internal/constant"
	"github.com/thoohv5/person/internal/provide/database"
	"github.com/thoohv5/person/internal/util"
)

// serverCmd represents the base command when called without any subcommands
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "数据库迁移工具",
	Long: `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 注册Config
		config, cleanConfig, err := boot.RegisterConfig(constant.ConfigPath(dir))
		if err != nil {
			return
		}
		// 注册Logger
		log, cleanLogger, err := boot.RegisterLogger(config)
		if err != nil {
			return
		}
		collection, err := im.GetCollection(util.Strikethrough2Underline(boot.Name))
		if err != nil {
			return
		}
		// 数据库创建
		db, ok := config.GetDatabase()[boot.Name]
		if !ok {
			return
		}
		_, cleanDB, err := database.New(db, log, database.WithCommand(args...), database.WithCollection(collection))
		if err != nil {
			return
		}
		cleanDB()
		cleanLogger()
		cleanConfig()
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}
