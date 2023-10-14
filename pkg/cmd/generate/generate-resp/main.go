package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/generate-resp/template"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
)

var tables = flag.String("t", "", "t: table, need register tables")
var project = flag.String("p", "ipam-interface", "p: project name")
var dir = flag.String("d", ".", "d: dir, Response path")

func main() {
	flag.Parse()
	if *tables == "" {
		fmt.Println("-t is required")
		return
	}
	for _, table := range strings.Split(*tables, ",") {
		baseParam, param := util.GenParam(*dir, *project, table, "")
		if err := util.GenFile("Response", baseParam, param, template.Response); err != nil {
			fmt.Println(err.Error())
			return
		}

	}
	fmt.Println("success")
}
