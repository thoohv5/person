package model

import (
	"reflect"

	"github.com/go-pg/pg/v10"

	ic "github.com/thoohv5/person/internal/context"
)

type Options struct {
	Db pg.DBI

	Models    []interface{}
	TableExpr func(pg.DBI, *pg.Query) *pg.Query
	QueryFunc func(*pg.Query)
	Result    *Result
	Msg       ic.ILogicalMessage

	updateZero bool
}

type Option func(o *Options)

func WithDb(db pg.DBI) Option {
	return func(o *Options) {
		o.Db = db
	}
}

func WithModel(model interface{}) Option {
	return func(o *Options) {
		if !reflect.ValueOf(model).IsNil() {
			o.Models = []interface{}{model}
		} else {
			o.Models = nil
		}
	}
}

func WithTableExpr(tableExpr func(pg.DBI, *pg.Query) *pg.Query) Option {
	return func(o *Options) {
		o.TableExpr = tableExpr
	}
}

func WithModels(models ...interface{}) Option {
	tempModels := make([]interface{}, 0, len(models))
	for _, model := range models {
		if !reflect.ValueOf(model).IsNil() {
			tempModels = append(tempModels, model)
		}
	}
	return func(o *Options) {
		o.Models = append(o.Models, tempModels...)
	}
}

func WithQueryFunc(queryFunc func(*pg.Query)) Option {
	return func(o *Options) {
		o.QueryFunc = queryFunc
	}
}

func WithUpdate(queryFunc func(*pg.Query)) Option {
	return func(o *Options) {
		o.QueryFunc = queryFunc
	}
}

func WithUpdateZero(updateZero bool) Option {
	return func(o *Options) {
		o.updateZero = updateZero
	}
}

type Result struct {
	// RowsAffected returns the number of rows affected by SELECT, INSERT, UPDATE,
	// or DELETE queries. It returns -1 if query can't possibly affect any rows,
	// e.g. in case of CREATE or SHOW queries.
	RowsAffected int32

	// RowsReturned returns the number of rows returned by the query.
	RowsReturned int32
}

func (r *Result) Gain(rowsAffected int, rowsReturned int) {
	r.RowsAffected = int32(rowsAffected)
	r.RowsReturned = int32(rowsReturned)
}

func WithResult(result *Result) Option {
	return func(o *Options) {
		o.Result = result
	}
}

func WithMessage(msg ic.ILogicalMessage) Option {
	return func(o *Options) {
		o.Msg = msg
	}
}
