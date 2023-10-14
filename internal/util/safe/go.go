package safe

import (
	"context"
	"runtime"

	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

// Go 使用此方法代替内部直接使用 go func
func Go(f func() error, log log.Logger) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]

				// 日志
				log.Errorc(context.Background(), "go func: panic recovered: ", logger.FieldInterface("err", r), logger.FieldInterface("buf", string(buf)))
			}
		}()
		if err := f(); err != nil {
			log.Errorc(context.Background(), "go run func err:", logger.FieldError(err))
		}
	}()
}

// GoWithCtx GO
func GoWithCtx(ctx context.Context, f func(ctx context.Context) error, log log.Logger) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]

				// 日志
				log.Errorc(ctx, "go func: panic recovered: ", logger.FieldInterface("err", r), logger.FieldInterface("buf", string(buf)))
			}
		}()
		if err := f(ic.Copy(ctx)); err != nil {
			log.Errorc(ctx, "go run func err:", logger.FieldError(err))
		}
	}()
}
