package boot

import (
	predis "github.com/thoohv5/person/internal/provide/redis"
	"github.com/thoohv5/person/pkg/log"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
)

// RegisterRedis 注册Redis
func RegisterRedis(conf config.Config, log log.Logger) (predis.IRedis, func(), error) {
	return predis.NewRedis(conf.GetRedis(), log)
}
