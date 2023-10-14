// Package nats .
package nats

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"

	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/transport"
)

// INatsConsumer Nats标准
type INatsConsumer = []transport.Server

// NewNatsConsumer 创建Nats消费者
func NewNatsConsumer(
	cfg *Config,
	log log.Logger,
	hs ...IHandle,
) INatsConsumer {
	items := INatsConsumer{}
	for i := 0; i < len(hs); i++ {
		if hs[i].Config() == nil {
			log.Warnc(context.Background(), fmt.Sprintf("[NATS] %v config invalid", reflect.TypeOf(hs[i]).Elem().Name()))
			continue
		}
		items = append(items, NewConsumer(log, hs[i], cfg))
	}
	return items
}

// IConsumer 消费者通用标准
type IConsumer transport.Server

// IHandle 实际消费者标准
type IHandle interface {
	// Config 配置
	Config() *Consumer
	// Handle 处理
	Handle(ctx context.Context, msg []byte) error
}

// consumer 消息的处理类
type consumer struct {
	IBase
	log      log.Logger
	h        IHandle
	sb       *nats.Subscription
	conf     *Config
	cfg      *Consumer
	stopChan chan bool
}

// NewConsumer 创建消费处理类
func NewConsumer(
	log log.Logger,
	h IHandle,
	baseCfg *Config,
) IConsumer {
	return &consumer{
		IBase:    NewBase(log, baseCfg),
		log:      log,
		h:        h,
		conf:     baseCfg,
		cfg:      h.Config(),
		stopChan: make(chan bool),
	}
}

// Start 启动
func (h *consumer) Start(ctx context.Context) error {
	if len(h.conf.GetUrl()) == 0 || len(h.cfg.GetSubject()) == 0 {
		return nil
	}
	// 创建链接
	if err := h.Conn(ctx); err != nil {
		return err
	}
	// 创建消费者
	if err := h.CreateConsumer(ctx, h.cfg); err != nil {
		return err
	}
	h.log.Infoc(ctx, "[NATS] consumer start", logger.FieldInterface("name", h.cfg.Name))
	// 订阅&处理
	if err := h.handle(ctx); err != nil {
		return err
	}
	return nil

}

// Stop 停止消费者
func (h *consumer) Stop(ctx context.Context) error {
	if len(h.conf.GetUrl()) == 0 || len(h.cfg.GetSubject()) == 0 {
		return nil
	}
	close(h.stopChan)
	// 关闭订阅
	if h.sb != nil && h.sb.IsValid() {
		if err := h.sb.Unsubscribe(); err != nil {
			h.log.Errorc(ctx, "[NATS] consumer unsubscribe error", logger.FieldError(err))
			return nil
		}
	}
	// 关闭链接
	if err := h.Close(ctx); err != nil {
		return err
	}
	h.log.Infoc(ctx, "[NATS] consumer stop", logger.FieldInterface("name", h.cfg.GetName()))
	return nil
}

func (h *consumer) handle(ctx context.Context) error {
	// 订阅
	sb, err := h.GetJetStream(ctx).PullSubscribe(h.cfg.GetSubject(), h.cfg.GetName(), nats.Bind(h.GetStreamName(ctx, h.cfg), h.cfg.GetName()))
	if err != nil {
		h.log.Errorc(ctx, "[NATS] PullSubscribe error", logger.FieldMap(
			map[string]interface{}{
				"name":    h.cfg.GetName(),
				"subject": h.cfg.GetSubject(),
				"err":     err,
			},
		))
		return err
	}
	h.sb = sb
	fetch := int32(200)
	if fh := h.cfg.GetFetch(); fh > 0 {
		fetch = fh
	}
	base := ic.FromContextBase(ctx)
	for {
		select {
		case <-h.stopChan:
			h.log.Infoc(ctx, "[NATS] consumer [for-fetch] stop", logger.FieldMap(
				map[string]interface{}{
					"name":    h.cfg.GetName(),
					"subject": h.cfg.GetSubject(),
				},
			))
			return nil
		default:
			// 拉取消息
			msgs := make([]*nats.Msg, 0, fetch)
			msgs, err = h.sb.Fetch(int(fetch))
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				} else if errors.Is(err, nats.ErrBadSubscription) { // 这里订阅被关闭了，所以需要退出
					return nil
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
			for i := 0; i < len(msgs); i++ {
				msg := msgs[i]
				msgCtx := context.Background()
				if values := msg.Header.Values(base.Key()); len(values) > 0 {
					if err = base.Unmarshal(values[0]); err != nil {
						h.log.Errorc(msgCtx, "[NATS] header err", logger.FieldError(err))
						return err
					}
					msgCtx, err = ic.WithMessage(msgCtx, base)
					if err != nil {
						h.log.Errorc(msgCtx, "[NATS] withMessage err", logger.FieldError(err))
						return err
					}
				}
				if err = func() (err error) {
					defer func() {
						if rev := recover(); rev != nil {
							buf := make([]byte, 64<<10)
							buf = buf[:runtime.Stack(buf, false)]
							err = fmt.Errorf("handlerr: panic recovered: %s\n%s", rev, buf)
						}
					}()
					if err = h.h.Handle(msgCtx, msg.Data); err != nil {
						return
					}
					return
				}(); err != nil {
					h.log.Errorc(msgCtx, "[NATS] handle err", logger.FieldError(err))
					continue
				}
				// ACK
				for defaultTick := 0; defaultTick < 3; defaultTick++ {
					time.Sleep(200 * time.Duration(defaultTick) * time.Millisecond)
					if err = msg.Ack(); err != nil {
						h.log.Errorc(msgCtx, "[NATS] ack err", logger.FieldError(err), logger.FieldInterface("times", defaultTick))
						continue
					}
					break
				}

			}
		}
	}
}
