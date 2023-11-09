package util

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	astutil2 "github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

type BaseParam struct {
	Table   string
	Project string
	Dir     string
	Package string
}

func GenParam(dir, project, table, tableName string) (*BaseParam, *TableWrapper) {
	baseParam := &BaseParam{
		Project: project,
		Dir:     dir,
		Table:   table,
		Package: strings.ReplaceAll(table, "_", ""),
	}
	upperCamelName := UpperCamelName(table)
	camelName := strings.ToLower(upperCamelName[:1]) + upperCamelName[1:]
	param := &TableWrapper{
		Name:          table,
		UpperCamlName: upperCamelName,
		CamlName:      camelName,
		MidlineName:   strings.ReplaceAll(table, "_", "-"),
		ProjectName:   project,
		Backquote:     "`",
		NameTrans:     tableName,
		Package:       strings.ReplaceAll(table, "_", ""),
	}
	return baseParam, param
}

func GenFile(kind string, baseParam *BaseParam, param *TableWrapper, format string) error {
	tpl, err := Execute(param, format)
	if err != nil {
		return errors.New(fmt.Sprintf("fail to generate %s[%s], err: %s\n", kind, baseParam.Table, err.Error()))

	}
	path := fmt.Sprintf("%s/%s.go", baseParam.Dir, baseParam.Table)
	_, err = os.Stat(path)
	if err == nil {
		return errors.New(fmt.Sprintf("fail to generate %s[%s], err: %s is exist\n", kind, baseParam.Table, path))

	}
	if err := ioutil.WriteFile(path, []byte(tpl), 0755); err != nil {
		return errors.New(fmt.Sprintf("fail to generate %s[%s], err: fail to write file. %s\n", kind, baseParam.Table, err.Error()))
	}
	return nil
}

// GenDao 生产dao，方便外部调用
func GenDao(dir, table string, tpl map[string]string) error {
	dirName := strings.ReplaceAll(table, "_", "")
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", dir, dirName), 0755); err != nil {
		return errors.New(fmt.Sprintf("dir create err %v\n", err))
	}
	_, param := GenParam(dir, "", table, "")
	for tag, str := range tpl {
		tpl, err := Execute(param, str)
		if err != nil {
			return errors.New(fmt.Sprintf("fail to generate dao: [%s].%s, err: %s\n", dirName, tag, err.Error()))
		}
		path := fmt.Sprintf("%s/%s/%s.go", dir, dirName, tag)
		_, err = os.Stat(path)
		if err == nil {
			return errors.New(fmt.Sprintf("fail to generate dao: [%s].%s, err: %s is exist\n", dirName, tag, path))
		}
		if err := ioutil.WriteFile(path, []byte(tpl), 0755); err != nil {
			return errors.New(fmt.Sprintf("fail to generate dao: [%s].%s, err: fail to write file. %s\n", dirName, tag, err.Error()))
		}
	}
	return nil
}

func RegController(fst *token.FileSet, f *ast.File, table string, regPath string) error {
	upperCamlName := UpperCamelName(table)

	// 往ProviderSet注入参数
	ret := astutil2.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\tNew%s", upperCamlName))
	if !ret {
		return errors.New(fmt.Sprintf("fail to rewrite ProviderSet"))
	}
	// astutil.PrintResult(fst,f)
	if err := astutil2.WriteToFile(fst, f, regPath); err != nil {
		return errors.New(fmt.Sprintf("fail to rewrite file"))
	}
	return nil
}

func RegService(fst *token.FileSet, f *ast.File, table string, regPath string) error {
	upperCamlName := UpperCamelName(table)

	// 往ProviderSet注入参数
	ret := astutil2.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\tNew%s", upperCamlName))
	if !ret {
		return errors.New(fmt.Sprintf("fail to rewrite ProviderSet"))
	}
	// astutil.PrintResult(fst,f)
	if err := astutil2.WriteToFile(fst, f, regPath); err != nil {
		return errors.New(fmt.Sprintf("fail to rewrite file"))
	}
	return nil
}

func RegRouter(fst *token.FileSet, f *ast.File, parentName, table, regPath string) error {
	_, param := GenParam("", "", table, "")
	// 往RegisterRouter方法添加参数
	ret := astutil2.AddParamToFunc(f, "RegisterRouter", fmt.Sprintf("\n\t%s", param.CamlName), fmt.Sprintf("*controller.%s", param.UpperCamlName))
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}
	// 新增参数
	varName := fmt.Sprintf("\n\tr%s", param.UpperCamlName)
	v := fmt.Sprintf("%s.Group(\"/%s\")", parentName, param.MidlineName)
	ret = astutil2.AddVarToFunc(f, "RegisterRouter", varName, v, "", "var")
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}
	// 新增代码块
	data := []astutil2.AstCallExpr{
		{
			FunName: varName,
			FunSel:  "POST",
			Args:    []string{astutil2.AddQuote(""), fmt.Sprintf("%s.Create", param.CamlName)},
		},
		{
			FunName: varName,
			FunSel:  "PUT",
			Args:    []string{astutil2.AddQuote("/:id"), fmt.Sprintf("%s.Update", param.CamlName)},
		},
		{
			FunName: varName,
			FunSel:  "GET",
			Args:    []string{astutil2.AddQuote(""), fmt.Sprintf("%s.List", param.CamlName)},
		},
		{
			FunName: varName,
			FunSel:  "GET",
			Args:    []string{astutil2.AddQuote("/all"), fmt.Sprintf("%s.All", param.CamlName)},
		},
		{
			FunName: varName,
			FunSel:  "GET",
			Args:    []string{astutil2.AddQuote("/:id"), fmt.Sprintf("%s.Detail", param.CamlName)},
		},
		{
			FunName: varName,
			FunSel:  "DELETE",
			Args:    []string{astutil2.AddQuote(""), fmt.Sprintf("%s.Delete", param.CamlName)},
		},
	}
	astutil2.AddCallBlockToFunc(f, "RegisterRouter", data, varName)
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}
	// astutil.PrintResult(fst,f)
	if err := astutil2.WriteToFile(fst, f, regPath); err != nil {
		return errors.New("fail to rewrite file")
	}
	return nil
}
