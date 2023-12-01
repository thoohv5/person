package db

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/pkg/log"
)

// tracingHook is a pg.QueryHook that adds OpenTelemetry instrumentation.
type loggerHook struct {
	log.Logger
}

func NewLoggerHook(logger log.Logger) pg.QueryHook {
	return &loggerHook{Logger: logger}
}

var _ pg.QueryHook = (*loggerHook)(nil)

func (h *loggerHook) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h *loggerHook) AfterQuery(ctx context.Context, evt *pg.QueryEvent) error {
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
