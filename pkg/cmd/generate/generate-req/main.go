package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/generate-req/template"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
)

var tables = flag.String("t", "", "t: table, need register tables")
var project = flag.String("p", "ipam-interface", "p: project name")
var dir = flag.String("d", ".", "d: dir, Request path")

func main() {
	flag.Parse()
	if *tables == "" {
		fmt.Println("-t is required")
		return
	}
	for _, table := range strings.Split(*tables, ",") {
		baseParam, param := util.GenParam(*dir, *project, table, "")
		if err := util.GenFile("Request", baseParam, param, template.Request); err != nil {
			fmt.Println(err.Error())
			return
		}

	}
	fmt.Println("success")
}
