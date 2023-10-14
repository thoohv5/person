package boot

import (
	"net/http"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	pHttp "github.com/thoohv5/person/internal/provide/http"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/transport"
)

// RegisterHTTP 注册HTTP
func RegisterHTTP(
	conf config.Config,
	log log.Logger,
	flag bool,
	handler http.Handler,
) transport.Server {
	return pHttp.New(conf.GetHttp(), log, flag, handler)
}
