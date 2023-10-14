package demo

import (
	"context"
	"time"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/internal/provide/nats"
	"github.com/thoohv5/person/pkg/log"
)

// 测试
type demo struct {
	log  log.Logger
	conf config.Config
}

// IDemo 测试
type IDemo nats.IHandle

// NewDemo 初始化消费者
func NewDemo(
	log log.Logger,
	conf config.Config,
) IDemo {
	sc := &demo{
		log:  log,
		conf: conf,
	}
	return sc
}

// Config 配置
func (d *demo) Config() *nats.Consumer {
	return d.conf.GetNats().Demo
}

// Handle 处理
func (d *demo) Handle(ctx context.Context, msg []byte) error {
	d.log.Infoc(ctx, "======================> ", logger.FieldString("time", time.Now().String()))
	return nil
}
