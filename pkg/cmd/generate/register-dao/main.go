package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/register-dao/ipam"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table, need register tables")
var project = flag.String("p", "ipam-interface", "p: project name")
var dir = flag.String("d", "./repository.go", "d: dir, Registered file path")

func main() {
	flag.Parse()

	path := *dir
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s not exist,err:%s\n", path, err.Error())
		return
	}
	fst, f := astutil.InitEnv(path)

	for _, table := range strings.Split(*tables, ",") {
		p, _ := util.GenParam(*dir, *project, table, "")
		switch *project {
		case "ipam-interface":
			if err := ipam.RegisterDao(fst, f, p, path); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

	fmt.Println("success")
}
