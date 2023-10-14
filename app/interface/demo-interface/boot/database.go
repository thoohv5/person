package boot

import (
	"errors"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/app/interface/demo-interface/migrations"
	"github.com/thoohv5/person/internal/provide/database"
	"github.com/thoohv5/person/internal/util"
	"github.com/thoohv5/person/pkg/log"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
)

// RegisterDatabase 注册DB
func RegisterDatabase(conf config.Config, log log.Logger) (pg.DBI, func(), error) {
	collection, err := migrations.GetCollection(util.Strikethrough2Underline(Name))
	if err != nil {
		return nil, nil, err
	}
	db, ok := conf.GetDatabase()[Name]
	if !ok {
		return nil, nil, errors.New("config not exists")
	}
	return database.New(db, log, database.WithCollection(collection))
}
