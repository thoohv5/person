package demo

import (
	"context"
	"time"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/internal/provide/cron"
	pcron "github.com/thoohv5/person/internal/provide/cron"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

// IDemo 测试
type IDemo pcron.ICron

type demo struct {
	cron.Base
	conf *cron.Timer
	log  log.Logger
}

// NewDemo 创建
func NewDemo(
	conf config.Config,
	log log.Logger,
) IDemo {
	return &demo{
		conf: conf.GetCron().Demo,
		log:  log,
	}
}

// Spec 获取Spec
func (d *demo) Spec() string {
	return d.conf.Spec
}

// Run 运行
func (d *demo) Run(ctx context.Context) func() {
	return func() {
		d.log.Infoc(ctx, "======================> ", logger.FieldString("time", time.Now().String()))
	}
}
