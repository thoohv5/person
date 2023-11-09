// Package boot 启动
package boot

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	RegisterConfig,
	RegisterLogger,
	RegisterDB,
	RegisterRedis,
	RegisterCron,
	RegisterNatsProducer,
	RegisterNatsConsumer,
	RegisterHTTP,
	InitApp,
)
