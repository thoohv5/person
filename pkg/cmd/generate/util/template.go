package util

import (
	"bytes"
	"text/template"
)

type TableWrapper struct {
	// 表名称
	Name string
	// 小驼峰
	CamlName string
	// 大驼峰
	UpperCamlName string
	// 中划线
	MidlineName string
	// 包名
	Package string
	// 项目名
	ProjectName string
	// 反引号
	Backquote string
	// 表名称解释
	NameTrans string
}

func Execute(param *TableWrapper, tpl string) (string, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("dao").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, param); err != nil {
		panic(err)
	}
	return buf.String(), nil
}
