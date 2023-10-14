package template

var Request = `package request

import (
	"github.com/thoohv5/person/internal/http"
)

// {{.UpperCamlName}}Create 创建时所需参数
type {{.UpperCamlName}}Create struct {
	
}

// {{.UpperCamlName}}Update 修改时所需参数
type {{.UpperCamlName}}Update struct {
	
}

// {{.UpperCamlName}}Delete 删除时所需参数
type {{.UpperCamlName}}Delete struct {
	IDList []int32
}

// {{.UpperCamlName}}All 获取全部数据时所需参数
type {{.UpperCamlName}}All struct {
	
}

// {{.UpperCamlName}}List 获取分页列表时所需参数
type {{.UpperCamlName}}List struct {
	http.BaseRequest
}

// {{.UpperCamlName}}PkID 获取详细数据时所需参数
type {{.UpperCamlName}}PkID struct {
	// ID
	ID int32 {{.Backquote}}uri:"id" binding:"required,number"{{.Backquote}}
}
`
