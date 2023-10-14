package http

type BaseRequest struct {
	// 页, 默认值 1
	PageNum int32 `form:"pageNum,default=1" json:"pageNum,default=1" default:"1"`
	// 页码，默认值 30
	PageSize int32 `form:"pageSize,default=30" json:"pageSize,default=30" default:"30"`
	// 排序字段
	SortField string `form:"sortField,default=id" json:"sortField,default=id" default:"id"`
	// 排序的方式，asc/desc
	SortOrder string `form:"sortOrder,default=asc" json:"sortOrder,default=asc" default:"asc"`
	// 排除Total
	ExclusiveTotal bool `json:"exclusiveTotal" example:"false"`
	// 排除List
	ExclusiveList bool `json:"exclusiveList" example:"false"`
	// Search 搜索条件
	Search []Search `form:"search" json:"search"`
}

// Search 搜索关键字
type Search struct {
	// Option 搜索选项 equal/like
	Option string `json:"option" binding:"required"`
	// Key 搜索键 多个参数逗号隔开表示或，比如a,b,c 搜索 a or b or c
	Key string `json:"key" binding:"required"`
	// Value 搜索值
	Value interface{} `json:"value" binding:"required"`
}

// ExtendAttrsParam 扩展属性
type ExtendAttrsParam struct {
	ExtendAttrs map[string]interface{} `json:"extendAttrs" binding:"omitempty"`
}
