package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table, Name of controller to be registered")
var dir = flag.String("d", "./controller.go", "d: dir, Registered file path")

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
		if err := util.RegController(fst, f, table, path); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	fmt.Println("success")
}
