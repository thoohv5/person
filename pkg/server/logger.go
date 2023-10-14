package server

import (
	"context"
	"fmt"
	"io"

	"github.com/thoohv5/person/pkg/log"
)

type devNullLog struct {
	w io.Writer
}

func NewDevNullLog(w io.Writer) log.Logger {
	return &devNullLog{
		w: w,
	}
}

func (d *devNullLog) AddCallerSkip(skip int) log.Logger {
	return d
}

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

func pre(ctx context.Context, fields ...log.Field) interface{} {
	fs := NewFields()
	// FieldContext(ctx)(fs)
	for _, field := range fields {
		field(fs)
	}
	return fs.Data()
}

func (d *devNullLog) Debugc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Infoc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Warnc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Errorc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
