// Package controller 控制器
package controller

import (
	"github.com/google/wire"

	"github.com/thoohv5/person/app/interface/demo-interface/internal/controller/demo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	demo.NewDemo,
)
