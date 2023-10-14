package constant

// ConfigPath 配置路径
type ConfigPath string

func (cf *ConfigPath) String() string {
	return string(*cf)
}

const (
	// SearchOptionLike 模糊搜索
	SearchOptionLike = "like"
	// SearchOptionStartWith 开头匹配
	SearchOptionStartWith = "rLike"
	// SearchOptionGT 大于
	SearchOptionGT = "gt"
	// SearchOptionGTE 大于等于
	SearchOptionGTE = "gte"
	// SearchOptionLT 小于
	SearchOptionLT = "lt"
	// SearchOptionLTE 小于等于
	SearchOptionLTE = "lte"
	// SearchOptionEqual 精确搜索
	SearchOptionEqual = "equal"
	// SearchOptionIn 在...中
	SearchOptionIn = "in"
	// SearchOptionBelong 属于
	SearchOptionBelong = "belong"
	// SearchOptionInclude 包含
	SearchOptionInclude = "include"
	// SearchOptionExist 存在
	SearchOptionExist = "exist"
)

// 排序方式
const (
	// SortAsc 升序
	SortAsc = "asc"
	// SortDesc 降序
	SortDesc = "desc"
)
