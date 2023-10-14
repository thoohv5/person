// Package request 请求
package request

import (
	"github.com/thoohv5/person/internal/provide/http"
)

// DemoCreate 创建时所需参数
type DemoCreate struct {
	// 名称
	Name string `json:"name" form:"name" binding:"required,normal"`
	// 年龄
	Age int32 `json:"age" form:"age" binding:"required,number"`
	// 状态
	State int32 `json:"state" form:"state" binding:"required,number,oneof=1 2"`
}

// DemoUpdate 修改时所需参数
type DemoUpdate struct {
	// 名称
	Name string `json:"name" form:"name" binding:"required,normal"`
	// 年龄
	Age int32 `json:"age" form:"age" binding:"required,number,gte=0"`
	// 状态
	State int32 `json:"state" form:"state" binding:"required,number,oneof=1 2"`
}

// DemoAll 获取全部数据时所需参数
type DemoAll struct {
	// 名称
	Name string `json:"name" form:"name" binding:"omitempty,normal"`
	// 年龄
	Age int32 `json:"age" form:"age" binding:"omitempty,number,gte=0"`
	// 状态
	State int32 `json:"state" form:"state" binding:"omitempty,oneof=1 2"`
}

// DemoList 获取分页列表时所需参数
type DemoList struct {
	http.BaseRequest
	// 名称
	Name string `json:"name" form:"name" binding:"omitempty,normal"`
	// 年龄
	Age int32 `json:"age" form:"age" binding:"omitempty,number,gte=0"`
	// 状态
	State int32 `json:"state" form:"state" binding:"omitempty,oneof=1 2"`
}

// DemoPkID 获取详细数据时所需参数
type DemoPkID struct {
	// ID
	ID int32 `uri:"id" binding:"required,number"`
}
