package demo

import (
	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
)

/**
等值 Field
IN  FieldIN
>   FieldGT
<   FieldLT
>=  FieldGTE
<=  FieldLTE
*/

// ID 条件
func ID(id int32) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("id = ?", id)
	}
}

// Name 条件
func Name(name string) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("name = ?", name)
	}
}

// State 条件
func State(state int32) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("state = ?", state)
	}
}
