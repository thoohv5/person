package ipam

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/thoohv5/person/pkg/cmd/generate/util"
	astutil "github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

func RegisterDao(fst *token.FileSet, f *ast.File, baseParam *util.BaseParam, registeredPath string) error {
	upperCamlName := util.CamelCase(baseParam.Table)
	camlName := strings.ToLower(upperCamlName[:1]) + upperCamlName[1:]
	// 添加import
	ret := astutil.AddImport(fst, f, camlName, fmt.Sprintf("github.com/thoohv5/person/app/interface/%s/internal/repository/%s", baseParam.Project, baseParam.Package))
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite import"))
	}
	// 往ProviderSet注入参数
	ret = astutil.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\t%s.New", camlName))
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite ProviderSet"))
	}

	// 往IRepository添加函数
	param := &astutil.AstFunc{
		Name: fmt.Sprintf("Get%s", upperCamlName),
		Results: []astutil.AstKv{
			{Key: "", Value: fmt.Sprintf("%s.I%s", camlName, upperCamlName)},
		},
	}
	ret = astutil.AddFuncToInterface(f, "IRepository", param)
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite IRepository"))
	}

	// 添加func
	param = &astutil.AstFunc{
		Name: fmt.Sprintf("Get%s", upperCamlName),
		Results: []astutil.AstKv{
			{Key: "", Value: fmt.Sprintf("%s.I%s", camlName, upperCamlName)},
		},
		Return: []string{fmt.Sprintf("d.%s", camlName)},
		Recv: &astutil.AstKv{
			Key:   "d",
			Value: "*repository",
		},
	}
	ret = astutil.AddFunc(f, param)

	// 往repository添加成员变量
	ret = astutil.AddKVToStruct(f, "repository", camlName, fmt.Sprintf("%s.I%s", camlName, upperCamlName), 12)
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite repository"))
	}

	// NewRepository函数添加参数和返回值
	ret = astutil.AddParamToFunc(f, "NewRepository", fmt.Sprintf("\n\t%s", camlName), fmt.Sprintf("%s.I%s", camlName, upperCamlName))
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite NewRepository"))
	}
	// 为NewRepository函数的变量rps添加数据
	ret = astutil.AddKVToFuncUnaryStruct(f, "NewRepository", "rps", fmt.Sprintf("\n\t%s", camlName), camlName)
	if !ret {
		return errors.New(fmt.Sprint("fail to rewrite variable rps of NewRepository"))
	}
	// astutil.PrintResult(fst,f)
	if err := astutil.WriteToFile(fst, f, registeredPath); err != nil {
		return errors.New(fmt.Sprint("fail to rewrite file"))
	}
	return nil
}
