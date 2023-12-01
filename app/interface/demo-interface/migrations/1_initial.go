// Package migrations 迁移
package migrations

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/repository/demo"
	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/util"
)

var (
	//go:embed *.sql
	sQLMigrations embed.FS
	models        = []interface{}{
		// 示例
		(*demo.Demo)(nil),
	}
)

// GetTables 获取表名
func GetTables() []string {
	return model.GetTables(models, sQLMigrations)
}

// GetCollection 获取collection
func GetCollection(tableName string) (*migrations.Collection, error) {
	tn := util.Strikethrough2Underline(tableName)
	collection := migrations.NewCollection().SetTableName(tn).DisableSQLAutodiscover(true)
	collection.MustRegisterTx(func(db migrations.DB) error {
		for _, m := range models {
			if err := db.Model(m).CreateTable(&orm.CreateTableOptions{
				IfNotExists: true,
			}); err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		for _, m := range models {
			if err := db.Model(m).DropTable(&orm.DropTableOptions{}); err != nil {
				return err
			}
		}
		return nil
	})
	if err := collection.DiscoverSQLMigrationsFromFilesystem(http.FS(sQLMigrations), "/"); err != nil {
		return nil, fmt.Errorf("sql migrations err: %w", err)
	}

	return collection, nil
}
