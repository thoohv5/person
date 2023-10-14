package boot

import (
	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	pcron "github.com/thoohv5/person/internal/provide/cron"
	"github.com/thoohv5/person/pkg/log"
)

// RegisterCron 注册日志
func RegisterCron(
	conf config.Config,
	log log.Logger,
	pcs []pcron.ICron,
) pcron.ICornServer {
	return pcron.NewCronServer(
		&conf.GetCron().Config,
		log,
		pcs...,
	)
}
