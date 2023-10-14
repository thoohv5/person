package model

type BaseRequest struct {
	// 数据开始位置
	Start int32
	// 返回数据条数
	Limit int32
	// 字段
	Fields []string
	// 排序：sort=["otc", "otc_type asc","created_at desc"]
	// 升序：没有或者asc标识，降序: desc标识
	Sorts []string
	// 分组
	GroupBy string
}
