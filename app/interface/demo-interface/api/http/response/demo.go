// Package response 返回值
package response

import "time"

// DemoEntity demo实体
type DemoEntity struct {
	// ID
	ID int32 `json:"id"`
	// 姓名
	Name string `json:"name"`
	// 年龄
	Age int32 `json:"age"`
	// 状态
	State int32 `json:"state"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
}
