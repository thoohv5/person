package main

import (
	"github.com/thoohv5/person/app/interface/demo-interface/api/docs"
	"github.com/thoohv5/person/app/interface/demo-interface/boot"
	"github.com/thoohv5/person/app/interface/demo-interface/cmd/command"
	"github.com/thoohv5/person/internal/util"
)

var version string

//	@title						Thooh demo API
//	@version					0.0.1
//	@description				this is thooh demo-interface project
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-User-Id
//	@BasePath					/
func main() {
	docs.SetVersion(util.Strikethrough2Underline(boot.Name), version)
	command.Execute(version)
}
