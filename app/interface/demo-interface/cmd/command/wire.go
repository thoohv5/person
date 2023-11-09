//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package command

import (
	"github.com/google/wire"

	"github.com/thoohv5/person/internal/constant"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/consumer"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/cron"
	"github.com/thoohv5/person/pkg"

	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/app/interface/demo-interface/boot"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/controller"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/repository"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/router"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/service"
	pcron "github.com/thoohv5/person/internal/provide/cron"
	"github.com/thoohv5/person/internal/provide/nats"
)

// initApp init
func initApp(dir constant.ConfigPath) (*pkg.App, func(), error) {
	panic(
		wire.Build(
			repository.ProviderSet,
			service.ProviderSet,
			cron.ProviderSet,
			consumer.ProviderSet,
			controller.ProviderSet,
			router.ProviderSet,
			boot.ProviderSet,
		),
	)
}

// initCron init
func initCron(dir constant.ConfigPath) ([]pcron.ICron, func(), error) {
	panic(
		wire.Build(
			boot.RegisterConfig,
			boot.RegisterLogger,
			// boot.RegisterDB,
			// repository.ProviderSet,
			// boot.RegisterRedis,
			// boot.RegisterNatsProducer,
			// service.ProviderSet,
			cron.ProviderSet,
		),
	)
}

type Module struct {
	ICS []pcron.ICron
	IHS []nats.IHandle
	IRS gin.RoutesInfo
}

func RegisterModule(
	ics []pcron.ICron,
	ihs []nats.IHandle,
	irs gin.RoutesInfo,
) *Module {
	return &Module{
		ICS: ics,
		IHS: ihs,
		IRS: irs,
	}
}

// InitProviderSet is init router providers.
var InitProviderSet = wire.NewSet(
	router.InitHTTP,
	wire.Bind(new(gin.IRouter), new(*gin.Engine)),
	router.RegisterRouter,
	ToRouter,
)

func ToRouter(
	engine *gin.Engine,
	flag bool,
) gin.RoutesInfo {
	if !flag {
		return nil
	}
	return engine.Routes()
}

// InitModule init
func InitModule(dir constant.ConfigPath, producer nats.IProducer) (*Module, func(), error) {
	panic(
		wire.Build(
			boot.RegisterConfig,
			boot.RegisterLogger,
			boot.RegisterDB,
			boot.RegisterRedis,
			repository.ProviderSet,
			service.ProviderSet,
			cron.ProviderSet,
			consumer.ProviderSet,
			controller.ProviderSet,
			InitProviderSet,
			RegisterModule,
		),
	)
}
