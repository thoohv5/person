package boot

import (
	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/internal/provide/nats"
	"github.com/thoohv5/person/pkg/log"
)

// RegisterNatsConsumer 注册Nats
func RegisterNatsConsumer(
	config config.Config,
	log log.Logger,
	ahs []nats.IHandle,
) nats.INatsConsumer {
	return nats.NewNatsConsumer(
		config.GetNats().Config,
		log,
		ahs...,
	)
}

// RegisterNatsProducer 注册NatsProducer
func RegisterNatsProducer(
	log log.Logger,
	config config.Config,
) nats.IProducer {
	a, err := nats.NewProducer(
		log,
		config.GetNats().Config,
	)
	if err != nil {
		panic(err)
	}
	return a
}
