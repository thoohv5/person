package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/register-field/template"
)

var dirs = flag.String("d", "", "d: dir, need register file dir")

func main() {
	flag.Parse()
	if *dirs == "" {
		fmt.Println("-d is required")
		return
	}
	mFieldDoc := make(map[string]string)
	for _, dir := range strings.Split(*dirs, ",") {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// 遍历文件列表并读取文件内容
		for _, file := range files {
			if !file.IsDir() {
				filePath := dir + "/" + file.Name()

				fset := token.NewFileSet()
				f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				sn := ""
				ast.Inspect(f, func(n ast.Node) bool {
					switch x := n.(type) {
					case *ast.TypeSpec:
						if _, ok := x.Type.(*ast.StructType); ok {
							sn = x.Name.Name
						}
					case *ast.Field:
						fn := ""
						if len(x.Names) > 0 {
							fn = x.Names[0].Name
						}
						if x.Doc != nil && len(fn) > 0 {
							ds := strings.Split(x.Doc.Text(), "\n")
							for i := 0; i < len(ds); i++ {
								if item := strings.TrimSpace(ds[i]); item != "" {
									mFieldDoc[fmt.Sprintf("%s.%s", sn, fn)] = item
									break
								}
							}
						}
					}
					return true
				})
			}
		}
		execute, err := template.Execute(&template.Wrapper{
			Package: "request",
			Field:   mFieldDoc,
		}, template.Tpl)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err = ioutil.WriteFile(fmt.Sprintf("%s/init.go", dir), []byte(execute), 0755); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
