package nats

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/thoohv5/person/pkg/log"
)

type base struct {
	log  log.Logger
	conf *Config
	nc   *nats.Conn
	js   nats.JetStreamContext
}

// IBase 标准
type IBase interface {
	// Conn 链接
	Conn(ctx context.Context) error
	// Close 关闭
	Close(ctx context.Context) error
	// GetJetStream 获取JetStream
	GetJetStream(ctx context.Context) nats.JetStreamContext
	// CreateConsumer 创建消费者
	CreateConsumer(ctx context.Context, consumerCfg *Consumer) error
	// GetStreamName 获取Stream名称
	GetStreamName(ctx context.Context, consumerCfg *Consumer) string
}

// NewBase 创建
func NewBase(
	log log.Logger,
	conf *Config,
) IBase {
	return &base{
		log:  log,
		conf: conf,
	}
}

// Conn 链接
func (b *base) Conn(ctx context.Context) error {
	// 链接
	conn, err := b.conn(ctx)
	if err != nil {
		return err
	}
	b.nc = conn

	// STREAM
	stream, err := b.jetStream(ctx)
	if err != nil {
		return err
	}
	b.js = stream

	// jetStream
	for _, streamConfig := range b.conf.GetStreams() {
		if err = b.createJetStream(ctx, streamConfig); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭
func (b *base) Close(_ context.Context) error {
	if b.nc != nil && !b.nc.IsClosed() {
		b.nc.Close()
	}
	return nil
}

// GetJetStream 获取JetStream
func (b *base) GetJetStream(ctx context.Context) nats.JetStreamContext {
	if b.js == nil {
		if err := b.Conn(ctx); err != nil {
			panic(err)
		}
	}
	return b.js
}

func (b *base) conn(ctx context.Context) (*nats.Conn, error) {
	ops := make([]nats.Option, 0)
	// 集群名称
	ops = append(ops, nats.Name(b.conf.GetClientName()))
	// 用户名密码
	if len(b.conf.GetUsername()) > 0 {
		ops = append(ops, nats.UserInfo(b.conf.GetUsername(), b.conf.GetPassword()))
	}
	// 重连间隔
	ops = append(ops, nats.ReconnectWait(time.Second*time.Duration(b.conf.GetReconnectTimeWait())))
	// 最大重连次数
	ops = append(ops, nats.MaxReconnects(int(b.conf.GetMaxReconnect())))
	// 重连回调
	ops = append(ops, nats.ReconnectHandler(func(nc *nats.Conn) {
		b.log.Warnc(ctx, fmt.Sprintf("[NATS] natscore.Reconnected [%s]", nc.ConnectedUrl()))
	}))
	// 发现服务器回调
	ops = append(ops, nats.DiscoveredServersHandler(func(nc *nats.Conn) {
		b.log.Infoc(ctx, fmt.Sprintf("[NATS] natscore.DiscoveredServersHandler %v", nc.DiscoveredServers()))
	}))
	// 服务器断开回调
	ops = append(ops, nats.ClosedHandler(func(nc *nats.Conn) {
		// b.log.Warnc(ctx, fmt.Sprintf("[NATS] natscore.ClosedHandler %v", nc.DiscoveredServers()))
	}))
	// 连接超时
	ops = append(ops, nats.Timeout(time.Duration(b.conf.GetConnectTimeout())*time.Second))
	// 连接失败回调
	ops = append(ops, nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] natscore.ErrorHandler %v", err))
	}))
	// 连接
	nc, err := nats.Connect(b.conf.GetUrl(), ops...)
	if err != nil {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] natscore.Connect %v", err))
		return nil, err
	}
	return nc, nil
}

func (b *base) jetStream(ctx context.Context) (nats.JetStreamContext, error) {
	// 获取stream
	js, err := b.nc.JetStream()
	if err != nil {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] nc JetStream err:%v", err))
		return nil, err
	}
	return js, nil
}

func (b *base) createJetStream(ctx context.Context, stream *JetStream) error {
	// 验证流是否存在
	_, err := b.js.StreamInfo(stream.Name)
	if err == nil {
		return nil
	}
	if !errors.Is(err, nats.ErrStreamNotFound) {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] Stream[%s] StreamInfo, err: %v", stream, err))
		return err
	}
	err = nil

	// 不存在，创建
	_, err = b.js.AddStream(b.initStreamConfig(stream))
	if err != nil {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] Stream[%s]  reate fail, error: %v", stream, err))
		return err
	}
	return nil
}

// CreateConsumer 创建消费者
func (b *base) CreateConsumer(ctx context.Context, consumerCfg *Consumer) error {
	// 获取stream名称
	stream := b.GetStreamName(ctx, consumerCfg)
	// 验证消费者是否存在
	_, err := b.js.ConsumerInfo(stream, consumerCfg.GetName())
	if err == nil {
		return nil
	}
	if !errors.Is(err, nats.ErrConsumerNotFound) {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] Consumer[%s] ConsumerInfo, err: %v", consumerCfg.GetName(), err))
		return err
	}
	err = nil
	// 初始化消费者配置
	cfg := b.initConsumerConfig(consumerCfg)
	// 不存在，创建
	_, err = b.js.AddConsumer(stream, cfg)
	if err != nil {
		b.log.Errorc(ctx, fmt.Sprintf("[NATS] Consumer[%s] create fail, error: %v", consumerCfg.GetName(), err))
		return err
	}
	return nil
}

// GetStreamName 获取Stream名称
func (b *base) GetStreamName(ctx context.Context, consumerCfg *Consumer) string {
	subject := consumerCfg.GetSubject()
	index := strings.Index(subject, ".")
	if index > -1 {
		return subject[:index]
	}
	b.log.Warnc(ctx, fmt.Sprintf("[NATS] Consumer[%s] Subject invalid", consumerCfg.GetName()))
	return subject
}

func (b *base) initStreamConfig(streamCfg *JetStream) *nats.StreamConfig {
	cfg := &nats.StreamConfig{
		Name:              streamCfg.Name,
		Subjects:          []string{fmt.Sprintf("%s.*", streamCfg.Name)},
		Storage:           nats.StorageType(streamCfg.Storage),
		Retention:         nats.LimitsPolicy,
		Discard:           nats.DiscardOld,
		MaxMsgs:           50000,
		MaxMsgsPerSubject: 10000,
		MaxAge:            24 * time.Hour,
		MaxBytes:          1 * 1024 * 1024 * 1024,
	}
	return cfg
}

func (b *base) initConsumerConfig(consumerCfg *Consumer) *nats.ConsumerConfig {
	cfg := &nats.ConsumerConfig{
		Durable:       consumerCfg.GetName(),
		DeliverPolicy: nats.DeliverNewPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	}
	if consumerCfg.GetDeliverPolicy() > 0 {
		cfg.DeliverPolicy = nats.DeliverPolicy(consumerCfg.GetDeliverPolicy() - 1)
	}
	if consumerCfg.GetAckPolicy() > 0 {
		cfg.AckPolicy = nats.AckPolicy(consumerCfg.GetAckPolicy() - 1)
	}
	return cfg
}
