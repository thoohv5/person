package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/generate-dao/template"
	"github.com/thoohv5/person/pkg/cmd/generate/register-dao/ipam"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table name, generate tables")
var cf = flag.String("c", "", "c: config file path：input absolute path. If empty, the default configuration will be used")
var project = flag.String("p", "ipam-interface", "p: project name")
var dir = flag.String("d", ".", "d: dir, service path")
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

	var registeredPath string
	var daoPath string
	for _, config := range cfg.Config {
		if config.Name == "repository" {
			registeredPath = fmt.Sprintf("%s/%s/%s", cfg.ProjectAbsPath, config.Path, config.RegisteredFile)
			daoPath = fmt.Sprintf("%s/%s", cfg.ProjectAbsPath, config.Path)
			break
		}
	}
	if registeredPath == "" {
		fmt.Printf("Registered file registeredPath not exit, registeredPath: %s\n", registeredPath)
		return
	}

	fst, f := astutil.InitEnv(registeredPath)
	for _, table := range strings.Split(*tables, ",") {
		// generate-dao
		baseTpl := map[string]string{table: template.ConfTemplate, "model": template.ConfTemplate3, "condition": template.ConfTemplate2}
		if err := util.GenDao(daoPath, table, baseTpl); err != nil {
			fmt.Printf("generate-dao: fail. err: %s\n", err.Error())
			return
		}
		fmt.Println("generate-dao: success")
		p, _ := util.GenParam(daoPath, cfg.ProjectName, table, "")
		// register-dap
		switch cfg.ProjectMode {
		case "ipam":
			if err := ipam.RegisterDao(fst, f, p, registeredPath); err != nil {
				fmt.Printf("register-dao: fail. err: %s\n", err.Error())
				return
			}
		default:
			fmt.Printf("Current project is not supported： %s\n", cfg.ProjectMode)
			return
		}
		fmt.Println("register-dao: success")
	}

	fmt.Println("init dao: success")

}
