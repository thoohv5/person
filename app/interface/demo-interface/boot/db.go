package boot

import (
	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/app/interface/demo-interface/migrations"
	"github.com/thoohv5/person/internal/provide/db"
	"github.com/thoohv5/person/pkg/log"
)

// RegisterDB 注册DB
func RegisterDB(
	conf config.Config,
	log log.Logger,
) (pg.DBI, func(), error) {
	return db.New(
		conf.GetDatabase(),
		log,
		db.WithApplicationName(Name),
		db.WithCollection(migrations.GetCollection),
	)
}
