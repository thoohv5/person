package log

import (
	"context"
	"reflect"
	"sort"
)

type IField interface {
	Set(key string, val interface{}, opts ...Option) IField
	Data() map[string]*Entity
}

func SMap(r map[string]*Entity, cb func(key string, value *Entity)) {
	ks := make([]string, 0, len(r))
	for k := range r {
		ks = append(ks, k)
	}
	// 对key排序
	sort.Strings(ks)
	for _, k := range ks {
		cb(k, r[k])
	}
}

type Field func(IField)

type Entity struct {
	Value interface{}
	Type  reflect.Kind
}

type Options struct {
	Type reflect.Kind
}
type Option func(*Options)

func WithType(oType reflect.Kind) Option {
	return func(o *Options) {
		o.Type = oType
	}
}

// Logger is a logger interface.
type Logger interface {
	Debugc(ctx context.Context, msg string, fields ...Field)
	Infoc(ctx context.Context, msg string, fields ...Field)
	Warnc(ctx context.Context, msg string, fields ...Field)
	Errorc(ctx context.Context, msg string, fields ...Field)
	AddCallerSkip(skip int) Logger
}
