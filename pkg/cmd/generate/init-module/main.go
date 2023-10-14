package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	tplctrl "github.com/thoohv5/person/pkg/cmd/generate/generate-ctrl/template"
	"github.com/thoohv5/person/pkg/cmd/generate/generate-dao/template"
	tplreq "github.com/thoohv5/person/pkg/cmd/generate/generate-req/template"
	tplresp "github.com/thoohv5/person/pkg/cmd/generate/generate-resp/template"
	tplserv "github.com/thoohv5/person/pkg/cmd/generate/generate-service/template"
	"github.com/thoohv5/person/pkg/cmd/generate/register-dao/ipam"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table name, generate tables")
var tableName = flag.String("tn", "{xxx}", "tn: table name")
var cf = flag.String("c", "", "c: config file path：input absolute path. If empty, the default configuration will be used")
var project = flag.String("p", "", "p: project name")
var dir = flag.String("d", "", "d: dir, project absolute path. If it is empty or '.', the configured default address will be used")
var mode = flag.String("m", "ipam", "m: project mode")

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
	var daoRegPath string
	var daoPath string
	var rRegPath string
	var rParentName string
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
		case "repository":
			daoRegPath = fmt.Sprintf("%s/%s/%s", cfg.ProjectAbsPath, config.Path, config.RegisteredFile)
			daoPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
		case "router":
			rRegPath = fmt.Sprintf("%s/%s/%s", cfg.ProjectAbsPath, config.Path, config.RegisteredFile)
			rParentName = config.ParentName
		}

	}

	fstDao, fDao := astutil.InitEnv(daoRegPath)
	fstCtrl, fCtrl := astutil.InitEnv(cRegPath)
	fstServ, fServ := astutil.InitEnv(sRegPath)
	fstRoute, fRoute := astutil.InitEnv(rRegPath)
	for _, table := range strings.Split(*tables, ",") {
		// gen-dao
		baseTpl := map[string]string{table: template.ConfTemplate, "model": template.ConfTemplate3, "condition": template.ConfTemplate2}
		if err := util.GenDao(daoPath, table, baseTpl); err != nil {
			fmt.Printf("generate-dao: fail. err: %s\n", err.Error())
			return
		}
		fmt.Println("generate-dao: success")

		// register-dap
		p, _ := util.GenParam(daoPath, cfg.ProjectName, table, *tableName)
		switch cfg.ProjectMode {
		case "ipam":
			if err := ipam.RegisterDao(fstDao, fDao, p, daoRegPath); err != nil {
				fmt.Printf("register-dao: fail. err: %s\n", err.Error())
				return
			}
		default:
			fmt.Printf("Current project is not supported： %s\n", cfg.ProjectName)
			return
		}
		fmt.Println("register-dao: success")

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

		// reg-router
		if err := util.RegRouter(fstRoute, fRoute, rParentName, table, rRegPath); err != nil {
			fmt.Printf("fail to register route, err:%s\n", err.Error())
			return
		}

	}
	fmt.Println("init logic: success")
}
