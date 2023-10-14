package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	tplctrl "github.com/thoohv5/person/pkg/cmd/generate/generate-ctrl/template"
	tplreq "github.com/thoohv5/person/pkg/cmd/generate/generate-req/template"
	tplresp "github.com/thoohv5/person/pkg/cmd/generate/generate-resp/template"
	tplserv "github.com/thoohv5/person/pkg/cmd/generate/generate-service/template"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table name, generate tables")
var cf = flag.String("c", "", "c: config file path：input absolute path. If empty, the default configuration will be used")
var project = flag.String("p", "ipam-interface", "p: project name")
var dir = flag.String("d", ".", "d: dir, service path")
var mode = flag.String("m", "ipam", "m: project mode")
var tableName = flag.String("tn", "{xxx}", "tn: table name")

//go:embed config.yaml
var defConfig string

func main() {
	flag.Parse()
	if *tables == "" {
		fmt.Println("-t is required")
		return
	}
	var cfg *util.Conf
	var err error
	if *cf == "" {
		cfg, err = util.InitConf2(defConfig)
		if err != nil {
			fmt.Printf("fail to init config, err: %s\n", err.Error())
			return
		}
	} else {
		cfg, err = util.InitConf(*cf)
		if err != nil {
			fmt.Printf("fail to init config, err: %s\n", err.Error())
			return
		}
	}

	if *project != "" {
		cfg.ProjectName = *project
	}

	if *mode != "" {
		cfg.ProjectMode = *mode
	}

	if *dir != "" && *dir != "." {
		cfg.ProjectAbsPath = *dir
	}
	fmt.Println("init config: success")

	var cPath string
	var cRegPath string
	var sPath string
	var sRegPath string
	var reqPath string
	var respPath string
	for _, config := range cfg.Config {
		switch config.Name {
		case "controller":
			cRegPath = fmt.Sprintf("%s/%s/%s", cfg.ProjectAbsPath, config.Path, config.RegisteredFile)
			cPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
		case "service":
			sRegPath = fmt.Sprintf("%s/%s/%s", cfg.ProjectAbsPath, config.Path, config.RegisteredFile)
			sPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
		case "request":
			reqPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
		case "response":
			respPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
		}

	}

	fstCtrl, fCtrl := astutil.InitEnv(cRegPath)
	fstServ, fServ := astutil.InitEnv(sRegPath)
	for _, table := range strings.Split(*tables, ",") {
		// gen-request
		baseParam, param := util.GenParam(reqPath, cfg.ProjectName, table, *tableName)
		if err := util.GenFile("Request", baseParam, param, tplreq.Request); err != nil {
			fmt.Printf("gen-request: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("gen-request: success")

		// gen-response
		baseParam, param = util.GenParam(respPath, cfg.ProjectName, table, *tableName)
		if err := util.GenFile("Response", baseParam, param, tplresp.Response); err != nil {
			fmt.Printf("gen-response: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("gen-response: success")

		// gen-controller
		baseParam, param = util.GenParam(cPath, cfg.ProjectName, table, *tableName)
		if err := util.GenFile("Controller", baseParam, param, tplctrl.Controller); err != nil {
			fmt.Printf("gen-controller: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("gen-controller: success")

		// reg-controller
		if err := util.RegController(fstCtrl, fCtrl, table, cRegPath); err != nil {
			fmt.Printf("reg-controller: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("reg-controller: success")

		// gen-service
		baseParam, param = util.GenParam(sPath, cfg.ProjectName, table, *tableName)
		if err := util.GenFile("Service", baseParam, param, tplserv.Service); err != nil {
			fmt.Printf("gen-service: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("gen-service: success")
		// reg-service
		if err := util.RegService(fstServ, fServ, table, sRegPath); err != nil {
			fmt.Printf("reg-service: false. err: %s\n", err.Error())
			return
		}
		fmt.Println("reg-service: success")

	}
	fmt.Println("init logic: success")
}
