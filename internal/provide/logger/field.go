package logger

import (
	"context"

	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/pkg/log"
)

type fields struct {
	data map[string]*log.Entity
}

func NewFields() log.IField {
	return &fields{
		data: make(map[string]*log.Entity),
	}
}

func (fs *fields) Set(key string, val interface{}, opts ...log.Option) log.IField {
	o := &log.Options{}
	for _, opt := range opts {
		opt(o)
	}
	fs.data[key] = &log.Entity{
		Value: val,
		Type:  o.Type,
	}
	return fs
}
func (fs *fields) Data() map[string]*log.Entity {
	return fs.data
}

func FieldInt32(key string, val int32) log.Field {
	return func(fs log.IField) {
		fs.Set(key, val)
	}
}

func FieldInt64(key string, val int64) log.Field {
	return func(fs log.IField) {
		fs.Set(key, val)
	}
}

func FieldString(key string, val string) log.Field {
	return func(fs log.IField) {
		fs.Set(key, val)
	}
}

func FieldMap(items map[string]interface{}) log.Field {
	return func(fs log.IField) {
		for k, v := range items {
			fs.Set(k, v)
		}
	}
}

func FieldContext(ctx context.Context) log.Field {
	return func(fs log.IField) {
		fs.Set("trace_id", ic.FromCtxTraceID(ctx))
	}
}

func FieldError(err error) log.Field {
	return func(fs log.IField) {
		if err == nil {
			return
		}
		fs.Set("err", err.Error())
	}
}

func FieldInterface(key string, val interface{}) log.Field {
	return func(fs log.IField) {
		fs.Set(key, val)
	}
}
