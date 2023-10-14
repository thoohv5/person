package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/generate-dao/template"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
)

var tables = flag.String("t", "", "t: table name, generate tables")
var showVersion = flag.Bool("version", false, "print the version and exit")
var dir = flag.String("d", ".", "d: generate file path：input absolute path")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("dao %v\n", release)
		return
	}
	if *tables == "" {
		fmt.Println("-t is required")
		return
	}
	for _, table := range strings.Split(*tables, ",") {
		baseTpl := map[string]string{table: template.ConfTemplate, "model": template.ConfTemplate3, "condition": template.ConfTemplate2}
		if err := util.GenDao(*dir, table, baseTpl); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	fmt.Println("success")
}
