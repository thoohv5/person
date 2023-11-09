package {{.Package}}

import "github.com/thoohv5/person/internal/validate"

// nolint:lll
var mFieldDoc = map[string]string{
{{range $key, $value := .Field -}}
"{{- $key}}": "{{ $value -}}",
{{- end}}
}

func init() {
	validate.RegisterFieldDoc(mFieldDoc)
}
