package boot

import (
	"fmt"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/internal/constant"
)

// RegisterConfig 注册Config
func RegisterConfig(dir constant.ConfigPath) (config.Config, func(), error) {
	cf, err := config.New(dir.String())
	if nil != err {
		return nil, nil, err
	}

	return cf, func() {
		if err = cf.Close(); err != nil {
			fmt.Printf("config Close err:%v\n", err)
			return
		}
	}, nil
}
