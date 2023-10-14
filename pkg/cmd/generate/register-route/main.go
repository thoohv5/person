package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var tables = flag.String("t", "", "t: table, need register tables")
var dir = flag.String("d", "./router.go", "d: dir, Registered file path")
var pName = flag.String("pn", "rIpam", "pn: parent name, upper-level route name")

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
		if err := util.RegRouter(fst, f, *pName, table, path); err != nil {
			fmt.Printf("fail to register route, err:%s\n", err.Error())
			return
		}
	}
	fmt.Println("success")

}
