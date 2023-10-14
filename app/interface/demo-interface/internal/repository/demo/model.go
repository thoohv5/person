// Package demo
// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package demo

import (
	"github.com/thoohv5/person/internal/model"
)

func init() {
	model.Register((*Demo)(nil))
}

// Demo 模块
type Demo struct {
	model.BaseModel
	tableName struct{} `pg:"demos,discard_unknown_columns"`

	Name  string `json:"name" pg:"name,notnull,unique,default:''"`
	Age   int32  `json:"age" pg:"age,notnull,default:0"`
	State int32  `json:"state" pg:"state,notnull,default:1"`
}
