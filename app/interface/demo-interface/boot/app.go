package boot

import (
	"github.com/thoohv5/person/internal/provide/cron"
	"github.com/thoohv5/person/internal/provide/nats"
	"github.com/thoohv5/person/pkg"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/transport"

	// docs
	_ "github.com/thoohv5/person/app/interface/demo-interface/api/docs"
	// translate
	_ "github.com/thoohv5/person/app/interface/demo-interface/api/translate"
)

const (
	// Name 服务标识
	Name = "demo-interface"
)

// InitApp 初始化App
func InitApp(
	logger log.Logger,
	hs transport.Server,
	cs cron.ICornServer,
	np nats.IProducer,
	as nats.INatsConsumer,
) *pkg.App {
	return pkg.New(
		pkg.Name(Name),
		pkg.Logger(logger),
		pkg.Server(
			append([]transport.Server{hs, cs, np}, as...)...,
		),
	)
}
