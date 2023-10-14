package boot

import (
	"github.com/thoohv5/person/pkg/log"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/internal/provide/logger"
)

// RegisterLogger 注册日志
func RegisterLogger(config config.Config) (log.Logger, func(), error) {
	return logger.New(config.GetLogger())
}
