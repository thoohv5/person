package main

import (
	_ "embed"
	"flag"
	"fmt"
	"go/token"
	"os"

	"github.com/pkg/errors"

	"github.com/thoohv5/person/pkg/cmd/generate/ctrl/template"
	"github.com/thoohv5/person/pkg/cmd/generate/util"
	"github.com/thoohv5/person/pkg/cmd/generate/util/astutil"
)

var dir = flag.String("d", "", "d: dir, root dir")
var projectPath = flag.String("path", "github.com/thoohv5/person", "path: project path")
var projectName = flag.String("pn", "demo-interface", "pn: project name")
var name = flag.String("n", "demo1", "n: name, controller name")
var remark = flag.String("r", "模版1", "r: remark, controller remark")

const (
	_path       = "app/interface"
	_request    = "api/http/request"
	_response   = "api/http/response"
	_controller = "internal/controller"
	_service    = "internal/service"
	_repository = "internal/repository"
	_router     = "internal/router"
)

func main() {
	flag.Parse()

	w := &template.CtrlWrapper{
		ProjectName: *projectName,
		Name:        *name,
		Remark:      *remark,
	}

	// 生成req
	if err := GenRequest(*dir, *name, w); err != nil {
		fmt.Println(err)
		return
	}

	// 生成resp
	if err := GenResponse(*dir, *name, w); err != nil {
		fmt.Println(err)
		return
	}

	// 生成Service
	if err := GenService(*dir, *name, w); err != nil {
		fmt.Println(err)
		return
	}

	// 注册Service
	if err := RegService(*dir, *name); err != nil {
		fmt.Println(err)
		return
	}

	// 生成Ctrl
	if err := GenController(*dir, *name, w); err != nil {
		fmt.Println(err)
		return
	}

	// 注册Ctrl
	if err := RegController(*dir, *projectPath, *projectName, *name); err != nil {
		fmt.Println(err)
		return
	}

	// 生成repository
	if err := GenRepository(*dir, *name, w); err != nil {
		fmt.Println(err)
		return
	}

	// 注册repository
	if err := RegRepository(*dir, *projectPath, *projectName, *name); err != nil {
		fmt.Println(err)
		return
	}

	// 注册router
	if err := RegRouter(*dir, *projectPath, *projectName, *name); err != nil {
		fmt.Println(err)
		return
	}

}

// GenRequest 生成
func GenRequest(dir string, name string, wrap *template.CtrlWrapper) error {
	path := fmt.Sprintf("%s/%s", dir, _request)
	content, err := template.Execute(wrap, template.GetRequest())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "fail to create dir")
		}
	}
	file := fmt.Sprintf("/%s/%s.go", path, name)
	if err = os.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(file); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}
	return nil
}

// GenResponse 生成
func GenResponse(dir string, name string, wrap *template.CtrlWrapper) error {
	path := fmt.Sprintf("%s/%s", dir, _response)
	content, err := template.Execute(wrap, template.GetResponse())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "fail to create dir")
		}
	}
	file := fmt.Sprintf("/%s/%s.go", path, name)
	if err = os.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(file); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}
	return nil
}

// GenService 生成
func GenService(dir string, name string, wrap *template.CtrlWrapper) error {
	path := fmt.Sprintf("%s/%s", dir, _service)
	content, err := template.Execute(wrap, template.GetService())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "fail to create dir")
		}
	}
	file := fmt.Sprintf("/%s/%s.go", path, name)
	if err = os.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(file); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}
	return nil
}

// RegService 注册
func RegService(dir string, name string) error {
	file := fmt.Sprintf("%s/%s/service.go", dir, _service)
	fst, f := astutil.InitEnv(file)
	// 往ProviderSet注入参数
	ret := astutil.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\tNew%s", util.UpperCamelName(name)))
	if !ret {
		return errors.New("fail to rewrite ProviderSet")
	}
	if err := astutil.WriteToFile(fst, f, file); err != nil {
		return errors.Wrap(err, "fail to rewrite file")
	}
	return nil
}

// GenController 生成
func GenController(dir string, name string, wrap *template.CtrlWrapper) error {
	path := fmt.Sprintf("%s/%s/%s", dir, _controller, name)
	content, err := template.Execute(wrap, template.GetController())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "fail to create dir")
		}
	}
	file := fmt.Sprintf("/%s/%s.go", path, name)
	if err = os.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(file); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}
	return nil
}

// RegController 注册
func RegController(dir string, projectPath, projectName, name string) error {
	file := fmt.Sprintf("%s/%s/controller.go", dir, _controller)
	fst, f := astutil.InitEnv(file)
	// 注入Import
	astutil.AddImport(fst, f, "", fmt.Sprintf("%s/%s/%s/%s/%s", projectPath, _path, projectName, _controller, name))
	// 往ProviderSet注入参数
	ret := astutil.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\t%s.New%s", name, util.UpperCamelName(name)))
	if !ret {
		return errors.New("fail to rewrite ProviderSet")
	}
	if err := astutil.WriteToFile(fst, f, file); err != nil {
		return errors.Wrap(err, "fail to rewrite file")
	}
	return nil
}

// GenRepository 生成
func GenRepository(dir string, name string, wrap *template.CtrlWrapper) error {
	path := fmt.Sprintf("%s/%s/%s", dir, _repository, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, "fail to create dir")
		}
	}

	mContent, err := template.Execute(wrap, template.GetModel())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	mFile := fmt.Sprintf("/%s/model.go", path)
	if err = os.WriteFile(mFile, []byte(mContent), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(mFile); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}

	cContent, err := template.Execute(wrap, template.GetCondition())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	cFile := fmt.Sprintf("/%s/condition.go", path)
	if err = os.WriteFile(cFile, []byte(cContent), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(cFile); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}

	dContent, err := template.Execute(wrap, template.GetDao())
	if err != nil {
		return errors.Wrap(err, "fail to gen template")
	}
	dFile := fmt.Sprintf("/%s/%s.go", path, name)
	if err = os.WriteFile(dFile, []byte(dContent), os.ModePerm); err != nil {
		return errors.Wrap(err, "fail to create file")
	}
	if _, err = os.Stat(dFile); err != nil {
		return errors.Wrap(err, "fail to check file not exist")
	}

	return nil
}

// RegRepository 注册
func RegRepository(dir string, projectPath, projectName, name string) error {
	file := fmt.Sprintf("%s/%s/repository.go", dir, _repository)
	fst, f := astutil.InitEnv(file)
	// 注入Import
	astutil.AddImport(fst, f, "", fmt.Sprintf("%s/%s/%s/%s/%s", projectPath, _path, projectName, _repository, name))
	// 往ProviderSet注入参数
	ret := astutil.AddValueToCaller(f, "ProviderSet", fmt.Sprintf("\n\t%s.New", name))
	if !ret {
		return errors.New("fail to rewrite ProviderSet")
	}
	// 往repository注入参数
	ret = astutil.AddKVToStruct(f, "repository", name, fmt.Sprintf("%v.I%v", name, util.UpperCamelName(name)), token.NoPos)
	if !ret {
		return errors.New("fail to rewrite repository")
	}
	// 往NewRepository注入参数
	ret = astutil.AddParamToFunc(f, "NewRepository", fmt.Sprintf("\n%v", name), fmt.Sprintf("%v.I%v", name, util.UpperCamelName(name)))
	if !ret {
		return errors.New("fail to rewrite NewRepository")
	}
	// 往NewRepository的rps注入数组
	ret = astutil.AddKVToFuncUnaryStruct(f, "NewRepository", "rps", fmt.Sprintf("\n%v", name), name)
	if !ret {
		return errors.New("fail to rewrite NewRepository rps")
	}
	// 往Repository注入func
	ret = astutil.AddFunc(f, &astutil.AstFunc{
		Recv: &astutil.AstKv{Key: "d", Value: " *repository"},
		Name: fmt.Sprintf("Get%v", util.UpperCamelName(name)),
		Results: []astutil.AstKv{
			{Value: fmt.Sprintf("%v.I%v", name, util.UpperCamelName(name))},
		},
		Return: []string{fmt.Sprintf("d.%v", name)},
	})
	if !ret {
		return errors.New("fail to rewrite fun")
	}
	if err := astutil.WriteToFile(fst, f, file); err != nil {
		return errors.Wrap(err, "fail to rewrite file")
	}
	return nil
}

// RegRouter 注册
func RegRouter(dir string, projectPath, projectName, name string) error {
	file := fmt.Sprintf("%s/%s/router.go", dir, _router)
	fst, f := astutil.InitEnv(file)
	// 注入Import
	astutil.AddImport(fst, f, "", fmt.Sprintf("%s/%s/%s/%s/%s", projectPath, _path, projectName, _controller, name))
	// 往ProviderSet注入参数
	ret := astutil.AddParamToFunc(f, "RegisterRouter", fmt.Sprintf("\n%v", name), fmt.Sprintf("*%s.%v", name, util.UpperCamelName(name)))
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}

	// 往ProviderSet注入数据
	// 新增参数
	varName := fmt.Sprintf("\n\tr%s", util.UpperCamelName(name))
	v := fmt.Sprintf("r.Group(\"/%s\")", util.Strikethrough(name))
	ret = astutil.AddVarToFunc(f, "RegisterRouter", varName, v, "", "var")
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}
	// 新增代码块
	data := []astutil.AstCallExpr{
		{
			FunName: varName,
			FunSel:  "POST",
			Args:    []string{astutil.AddQuote(""), fmt.Sprintf("%s.Create", name)},
		},
		{
			FunName: varName,
			FunSel:  "PUT",
			Args:    []string{astutil.AddQuote("/:id"), fmt.Sprintf("%s.Update", name)},
		},
		{
			FunName: varName,
			FunSel:  "GET",
			Args:    []string{astutil.AddQuote(""), fmt.Sprintf("%s.List", name)},
		},
		{
			FunName: varName,
			FunSel:  "GET",
			Args:    []string{astutil.AddQuote("/:id"), fmt.Sprintf("%s.Detail", name)},
		},
		{
			FunName: varName,
			FunSel:  "DELETE",
			Args:    []string{astutil.AddQuote(""), fmt.Sprintf("%s.Delete", name)},
		},
	}
	ret = astutil.AddCallBlockToFunc(f, "RegisterRouter", data, varName)
	if !ret {
		return errors.New("fail to rewrite RegisterRouter")
	}

	if err := astutil.WriteToFile(fst, f, file); err != nil {
		return errors.Wrap(err, "fail to rewrite file")
	}
	return nil
}
