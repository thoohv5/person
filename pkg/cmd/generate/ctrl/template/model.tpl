// Package {{ .Name }}
// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package {{ .Name }}

import (
	"github.com/thoohv5/person/internal/model"
)

func init() {
	model.Register((*{{ U .Name }})(nil))
}

// {{ .Remark }} 模块
type {{ U .Name }} struct {
	model.BaseModel
	tableName struct{} `pg:"{{ .Name }},discard_unknown_columns"`

}
