package cron

import (
	"context"
	"fmt"
	"runtime"

	"github.com/robfig/cron/v3"

	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/transport"
)

// ICornServer 定时标准
type ICornServer interface {
	transport.Server
}

// ICron 任务标准
type ICron interface {
	// Spec 任务计划
	Spec() string
	// Run 任务
	Run(ctx context.Context) func()
}

type cornServer struct {
	conf *Config
	log  log.Logger
	jobs []ICron
	c    *cron.Cron
}

// NewCronServer 创建 定时任务
func NewCronServer(
	conf *Config,
	log log.Logger,
	jobs ...ICron,
) ICornServer {
	opts := []cron.Option{
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(cron.VerbosePrintfLogger(NewDefaultCronLogger(log, WithIsError(true))))),
	}
	if conf.GetDebug() {
		opts = append(opts, cron.WithLogger(cron.VerbosePrintfLogger(NewDefaultCronLogger(log))))
	}
	return &cornServer{
		conf: conf,
		log:  log,
		jobs: jobs,
		c:    cron.New(opts...),
	}
}

// Start 开始
func (s *cornServer) Start(ctx context.Context) (err error) {
	if !s.conf.GetEnable() {
		return
	}

	for i := 0; i < len(s.jobs); i++ {
		item := s.jobs[i]
		if len(item.Spec()) == 0 {
			continue
		}
		if _, err = s.c.AddFunc(item.Spec(), func() {
			defer func() {
				if rev := recover(); nil != rev {
					buf := make([]byte, 64<<10)
					buf = buf[:runtime.Stack(buf, false)]
					err = fmt.Errorf("Corn: panic recovered: %s\n%s", rev, buf)
				}
			}()
			ctx = ic.Copy(ctx)
			item.Run(ctx)()
		}); err != nil {
			s.log.Errorc(ctx, "[CRON] server run err", logger.FieldError(err), logger.FieldInterface("item", item))
			return err
		}
	}

	s.c.Start()
	s.log.Infoc(ctx, "[CRON] server start")
	return nil
}

// Stop 停止
func (s *cornServer) Stop(ctx context.Context) (err error) {
	if !s.conf.GetEnable() {
		return
	}
	s.c.Stop()
	s.log.Infoc(ctx, "[CRON] server stop")
	return
}

type defaultCronLogger struct {
	log.Logger

	o *options
}

type options struct {
	ctx     context.Context
	isError bool
}

type Option func(os *options)

func WithIsError(isError bool) Option {
	return func(o *options) {
		o.isError = isError
	}
}

// NewDefaultCronLogger 创建定时日志
func NewDefaultCronLogger(logger log.Logger, opts ...Option) interface{ Printf(string, ...interface{}) } {
	o := &options{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(o)
	}

	return &defaultCronLogger{
		o:      o,
		Logger: logger,
	}
}

// Printf 日志输出
func (d *defaultCronLogger) Printf(msg string, data ...interface{}) {
	if !d.o.isError {
		d.Infoc(d.o.ctx, fmt.Sprintf(msg, data...))
		return
	}
	d.Errorc(d.o.ctx, fmt.Sprintf(msg, data...))
}
