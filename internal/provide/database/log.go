package database

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/thoohv5/person/pkg/log"
)

type queryOperation interface {
	Operation() orm.QueryOp
}

// LoggerHook is a pg.QueryHook that adds OpenTelemetry instrumentation.
type LoggerHook struct {
	log.Logger
}

func NewTracingHook() *LoggerHook {
	return new(LoggerHook)
}

var _ pg.QueryHook = (*LoggerHook)(nil)

func (h *LoggerHook) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h *LoggerHook) AfterQuery(ctx context.Context, evt *pg.QueryEvent) error {

	formattedQuery, err := evt.FormattedQuery()
	if err != nil {
		return err
	}
	query := string(formattedQuery)

	var (
		numRow int
	)
	fn, file, line := funcFileLine("github.com/go-pg/pg")

	if evt.Err != nil {
		switch evt.Err {
		case pg.ErrNoRows, pg.ErrMultiRows:
		default:
			err = evt.Err
		}
	} else if evt.Result != nil {
		numRow = evt.Result.RowsAffected()
		if numRow == 0 {
			numRow = evt.Result.RowsReturned()
		}
	}

	for _, keyword := range []string{"comptroller", "COMMIT", "BEGIN", "pg_logical_emit_message"} {
		if strings.Contains(query, keyword) {
			return nil
		}
	}

	h.Logger.Debugc(ctx, "database", func(field log.IField) {
		field.Set("filepath", file)
		field.Set("slineno", line)
		field.Set("function", fn)
		field.Set("statement", query, log.WithType(reflect.String))
		field.Set("err", err)
		field.Set("affect row", numRow)
	})

	return nil
}

func funcFileLine(pkg string) (string, string, int) {
	const depth = 16
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	ff := runtime.CallersFrames(pcs[:n])

	var fn, file string
	var line int
	for {
		f, ok := ff.Next()
		if !ok {
			break
		}
		fn, file, line = f.Function, f.File, f.Line
		if !strings.Contains(fn, pkg) {
			break
		}
	}

	if ind := strings.LastIndexByte(fn, '/'); ind != -1 {
		fn = fn[ind+1:]
	}

	return fn, file, line
}
