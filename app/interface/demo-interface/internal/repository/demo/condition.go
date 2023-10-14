package demo

import (
	"fmt"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/provide/http"
)

// ID 条件
func ID(id int32) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("id = ?", id)
	}
}

// GetCommonQuery 获取公共查询条件
func GetCommonQuery(params http.BaseRequest) model.QueryOption {
	demoPrefix := model.GetTableName((*Demo)(nil))
	return func(query *pg.Query) {
		model.GetCommonQuery(params, model.DefaultDeal(demoPrefix))(query)
		for _, s := range params.Search {
			// 特殊字段处理
			switch s.Key {
			default:
				s.Key = fmt.Sprintf("%s.%s", demoPrefix, s.Key)
				// 通用字段处理
				model.DealSearch(query, s)
			}
		}
	}
}

// Name 条件
func Name(name string) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("name = ?", name)
	}
}
