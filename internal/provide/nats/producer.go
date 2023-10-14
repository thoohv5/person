// Package nats .
package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"

	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

// IProducer 生产者标准
type IProducer interface {
	Start(ctx context.Context) error
	Publish(ctx context.Context, config *Producer, msg IMessage) error
	Stop(ctx context.Context) error
}

// producer 生产者
type producer struct {
	IBase

	conf *Config
	log  log.Logger
}

// NewProducer 创建生产者
func NewProducer(log log.Logger, conf *Config) (IProducer, error) {
	return &producer{
		IBase: NewBase(log, conf),

		conf: conf,
		log:  log,
	}, nil
}

// Start 启动生产者
func (p *producer) Start(ctx context.Context) error {
	if len(p.conf.GetUrl()) == 0 {
		return nil
	}
	if err := p.Conn(ctx); err != nil {
		return err
	}
	p.log.Infoc(ctx, "[NATS] producer start")
	return nil
}

// Publish 发布消息
func (p *producer) Publish(ctx context.Context, config *Producer, msg IMessage) error {
	if len(p.conf.GetUrl()) == 0 {
		return nil
	}
	if config == nil {
		p.log.Warnc(ctx, "[NATS] producer publish config invalid", logger.FieldInterface("config", config))
		return nil
	}
	// 序列化
	marshal, err := msg.Marshal()
	if err != nil {
		return err
	}
	// 获取ctx
	header := nats.Header{}
	lm := ic.FromContextBase(ctx)
	s, err := lm.Marshal(ctx)
	if err != nil {
		return err
	}
	header.Set(lm.Key(), s)

	// 循环3次
	sendMsg := &nats.Msg{
		Subject: config.GetSubject(),
		Header:  header,
		Data:    marshal,
	}
	for i := 0; i < 3; i++ {
		// 发送
		ack, err := p.GetJetStream(ctx).PublishMsg(sendMsg)
		if err != nil {
			p.log.Errorc(ctx, "[NATS] producer publish fail", logger.FieldError(err), logger.FieldMap(map[string]interface{}{
				"msg": sendMsg,
				"ack": ack,
				"err": err,
			}))
			time.Sleep(200 * time.Duration(i) * time.Millisecond)
			continue
		}
		break
	}
	// 发送失败，仅仅记录日志
	return nil
}

// Stop 停止生产者
func (p *producer) Stop(ctx context.Context) error {
	if len(p.conf.GetUrl()) == 0 {
		return nil
	}
	// 链接关闭
	if err := p.Close(ctx); err != nil {
		return err
	}
	p.log.Infoc(ctx, "[NATS] producer stop")
	return nil
}
