package template

import (
	"bytes"
	"text/template"

	"github.com/thoohv5/person/pkg/cmd/generate/util"
)

type CtrlWrapper struct {
	// 项目名称
	ProjectName string
	// 名称
	Name string
	// 备注
	Remark string
}

func Execute(param *CtrlWrapper, tpl string) (string, error) {
	buf := new(bytes.Buffer)
	temp, err := template.New("ctrl").Funcs(template.FuncMap{
		"U": util.UpperCamelName,
		"L": util.LowerCameName,
		"S": util.Strikethrough,
	}).Parse(tpl)
	if err != nil {
		return "", err
	}
	if err = temp.Execute(buf, param); err != nil {
		return "", err
	}
	return buf.String(), nil
}
