package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/extra/pgotel/v10"
	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

const defaultDBConnMaxRetries = 5

func New(dbConfig *Config, log log.Logger, opts ...Option) (pg.DBI, func(), error) {

	o := &options{
		command:    []string{"init", "version", "", "version"},
		collection: migrations.DefaultCollection,
	}
	for _, opt := range opts {
		opt(o)
	}

	opt, err := pg.ParseURL(fmt.Sprintf("%s://%s", dbConfig.GetDriver(), dbConfig.GetSource()))
	if err != nil {
		log.Errorc(context.Background(), "db Connect err", logger.FieldError(err), logger.FieldInterface("config", dbConfig))
		return nil, nil, err
	}

	// 最大链接时间
	opt.MaxConnAge = time.Second * time.Duration(dbConfig.GetConnMaxLifetimeSeconds())
	// 最小的空闲链接数
	opt.MinIdleConns = int(dbConfig.GetMinIdleConns())
	// 链接池的最大连接数
	opt.PoolSize = int(dbConfig.GetMaxOpenConns())
	// 最大retry次数
	if maxRetries := int(dbConfig.GetMaxRetries()); maxRetries == 0 {
		opt.MaxRetries = defaultDBConnMaxRetries
	} else {
		opt.MaxRetries = maxRetries
	}

	db := pg.Connect(opt)
	db.AddQueryHook(pgotel.NewTracingHook())
	db.AddQueryHook(&LoggerHook{log})
	// db.AddQueryHook(pgdebug.NewDebugHook())

	for _, c := range o.command {
		// 导入
		params := make([]string, 0)
		if len(c) > 0 {
			params = append(params, c)
		}
		ov, nv, err := o.collection.Run(db, params...)
		if err != nil {
			log.Debugc(context.Background(), "db Migrate err", logger.FieldError(err))
			return nil, nil, err
		}
		log.Debugc(context.Background(), "db Migrate info", logger.FieldMap(map[string]interface{}{
			"command": c,
			"old":     ov,
			"new":     nv,
		}))
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return db, func() {
		if err := db.Close(); nil != err {
			log.Errorc(context.Background(), "db Close err", logger.FieldError(err))
		}
	}, nil

}

type options struct {
	command    []string
	collection *migrations.Collection
}

type Option func(*options)

// WithCommand 设置migrations的command
func WithCommand(cs ...string) Option {
	return func(o *options) {
		o.command = cs
	}
}

// WithCollection 设置migrations的collection
func WithCollection(c *migrations.Collection) Option {
	return func(o *options) {
		o.collection = c
	}
}
