// Package consumer 消费者
package consumer

import (
	"github.com/google/wire"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/consumer/demo"
	"github.com/thoohv5/person/internal/provide/nats"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	RegisterNatsConsumer,
	demo.NewDemo,
)

// RegisterNatsConsumer 注册消费者
func RegisterNatsConsumer(
	demo demo.IDemo,
) []nats.IHandle {
	return []nats.IHandle{
		demo,
	}
}
