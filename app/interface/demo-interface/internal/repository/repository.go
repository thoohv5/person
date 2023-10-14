// Package repository 资源
package repository

import (
	"context"
	"errors"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	ic "github.com/thoohv5/person/internal/context"
	predis "github.com/thoohv5/person/internal/provide/redis"

	"github.com/go-pg/pg/v10"
	"github.com/google/wire"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/repository/demo"
	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/pkg/log"
)

// ProviderSet is repository providers.
var ProviderSet = wire.NewSet(
	NewRepository,
	demo.New,
)

// IRepository 资源标准
type IRepository interface {
	RunInTransaction(ctx context.Context, tx ...func(ctx context.Context, re IRepository, opts ...model.Option) error) error
	Begin(opts ...model.Option) (bool, []model.Option, error)
	Commit(opts ...model.Option) error
	Rollback(opts ...model.Option) error

	GetCache() predis.IRedis
	GetDemo() demo.IDemo
}

type repository struct {
	conf   config.Config
	logger log.Logger
	gdb    pg.DBI
	rdb    predis.IRedis

	// 注册 dao
	demo demo.IDemo
}

// NewRepository .
func NewRepository(
	conf config.Config,
	log log.Logger,
	gdb pg.DBI,
	rdb predis.IRedis,

	demo demo.IDemo,

) (IRepository, func(), error) {
	cleanup := func() {
	}

	rps := &repository{
		conf:   conf,
		logger: log,
		gdb:    gdb,
		rdb:    rdb,

		demo: demo,
	}

	return rps, cleanup, nil
}

func (d *repository) RunInTransaction(ctx context.Context, txs ...func(ctx context.Context, re IRepository, opts ...model.Option) error) error {
	//nolint:contextcheck
	ctx = ic.WithLM(ctx)
	return d.gdb.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for _, t := range txs {
			if err := t(ctx, d, model.WithDb(tx)); err != nil {
				return err
			}
		}
		if err := model.Message(ctx, tx); err != nil {
			return err
		}

		return nil
	})
}

func (d *repository) GetCache() predis.IRedis {
	return d.rdb
}

func (d *repository) Begin(opts ...model.Option) (bool, []model.Option, error) {
	o := &model.Options{}
	for _, op := range opts {
		op(o)
	}

	if _, ok := o.Db.(*pg.Tx); ok {
		return false, opts, nil
	}

	tx, err := d.gdb.Begin()
	if err != nil {
		return false, nil, err
	}
	return true, append([]model.Option{model.WithDb(tx)}, opts...), nil
}

func (d *repository) Commit(opts ...model.Option) error {
	o := &model.Options{}
	for _, op := range opts {
		op(o)
	}

	tx, ok := o.Db.(*pg.Tx)
	if !ok {
		return errors.New("invalid tx")
	}

	return tx.Commit()
}

func (d *repository) Rollback(opts ...model.Option) error {
	o := &model.Options{}
	for _, op := range opts {
		op(o)
	}

	tx, ok := o.Db.(*pg.Tx)
	if !ok {
		return errors.New("invalid tx")
	}

	return tx.Rollback()
}

func (d *repository) GetDemo() demo.IDemo {
	return d.demo
}
