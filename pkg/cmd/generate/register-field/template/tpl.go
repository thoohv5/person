package template

import (
	"bytes"
	"text/template"
)

type Wrapper struct {
	Package string
	Field   map[string]string
}

func Execute(param *Wrapper, tpl string) (string, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("field").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err = tmpl.Execute(buf, param); err != nil {
		panic(err)
	}
	return buf.String(), nil
}
