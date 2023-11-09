// Package router 路由
package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/controller/demo"
	"github.com/thoohv5/person/internal/util"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/thoohv5/person/app/interface/demo-interface/boot"
	ic "github.com/thoohv5/person/internal/context"
	pHttp "github.com/thoohv5/person/internal/provide/http"
	"github.com/thoohv5/person/internal/validate"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/server/middleware"
	"github.com/thoohv5/person/pkg/server/middleware/requestid"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(
	InitHTTP,
	wire.Bind(new(gin.IRouter), new(*gin.Engine)),
	wire.Bind(new(http.Handler), new(*gin.Engine)),
	RegisterRouter,
)

// InitHTTP 注册Router
func InitHTTP(
	conf config.Config,
	log log.Logger,
) (*gin.Engine, error) {
	if model := conf.GetHttp().GetModel(); model == gin.ReleaseMode {
		gin.SetMode(model)
	}
	engine := gin.New()
	engine.ContextWithFallback = true
	engine.UseRawPath = true

	engine.NoRoute(func(ctx *gin.Context) {
		pHttp.Fail(ctx, errors.New("404 Not Found "+ctx.Request.Method+" "+ctx.FullPath()))
	})
	engine.NoMethod(func(ctx *gin.Context) {
		pHttp.Fail(ctx, errors.New("404 Not Method "+ctx.Request.Method+" "+ctx.FullPath()))
	})

	hConfig := conf.GetHttp()

	if hConfig.GetEnablePprof() {
		pprof.Register(engine, fmt.Sprintf("/%s%s", boot.Name, pprof.DefaultPrefix))
	}

	engine.Use(
		middleware.Recovery(log, true),
		requestid.New(
			requestid.WithCustomHeaderStrKey(requestid.HeaderStrKey(ic.GetTraceIDLabel())),
		),
		middleware.Params(ic.GetLangLabel(), ic.GetUserIDLabel()),
		middleware.Transform(ic.WithHeader()),
		middleware.Logger(log, "2006-01-02T15:04:05Z08:00", false),
	)

	if err := validate.InitValidate("zh"); err != nil {
		return engine, err
	}

	if hConfig.GetEnableSwag() {
		relativePath := fmt.Sprintf("/%s/%s/*any", boot.Name, util.Strikethrough(boot.Name))
		fmt.Printf("swagger url: %s%s\n", conf.GetHttp().GetAddr(), relativePath)
		engine.GET(relativePath, ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1), func(config *ginSwagger.Config) {
			config.InstanceName = util.Strikethrough(boot.Name)
		}))
	}

	return engine, nil
}

// RegisterRouter 注册Router
func RegisterRouter(
	root gin.IRouter,
	demo *demo.Demo,
) bool {
	r := root.Group(fmt.Sprintf("/%s", boot.Name))

	rDemo := r.Group("/demo")
	{
		rDemo.POST("", demo.Create)
		rDemo.PUT("/:id", demo.Update)
		rDemo.GET("", demo.List)
		rDemo.GET("/:id", demo.Detail)
		rDemo.DELETE("/:id", demo.Delete)
	}
	return true
}
