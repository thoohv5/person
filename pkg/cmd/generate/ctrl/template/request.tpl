// Package request 请求
package request

import (
	"github.com/thoohv5/person/internal/provide/http"
)

// {{ U .Name }}Create {{ .Remark }}创建
type {{ U .Name }}Create struct {
}

// {{ U .Name }}Update {{ .Remark }}修改
type {{ U .Name }}Update struct {
}

// {{ U .Name }}List {{ .Remark }}列表
type {{ U .Name }}List struct {
	http.BaseRequest
}

// {{ U .Name }}PkID {{ .Remark }}详情
type {{ U .Name }}PkID struct {
}
