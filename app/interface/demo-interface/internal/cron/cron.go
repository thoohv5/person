// Package cron 定时任务
package cron

import (
	"github.com/google/wire"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/cron/demo"
	pcron "github.com/thoohv5/person/internal/provide/cron"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	RegisterCron,
	demo.NewDemo,
)

// RegisterCron 注册定时任务
func RegisterCron(
	demo demo.IDemo,
) []pcron.ICron {
	return []pcron.ICron{
		demo,
	}
}
