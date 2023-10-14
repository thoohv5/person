package template

var ConfTemplate = `// package {{.Package}}
// nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import (
	"context"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/http"
)

type {{.CamlName}} struct {
	db pg.DBI
}

// I{{.UpperCamlName}} 标准
type I{{.UpperCamlName}} interface {
	Create(ctx context.Context, info *{{.UpperCamlName}}, opts ...model.Option) (err error)
	Update(ctx context.Context, con model.IQuery, info *{{.UpperCamlName}}, opts ...model.Option) (err error)
	Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error)
	Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *{{.UpperCamlName}}, err error)
	List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*{{.UpperCamlName}}, err error)
	Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error)
}

// New 创建
func New(db pg.DBI) (I{{.UpperCamlName}}, error) {
	return &{{.CamlName}}{
		db: db,
	}, nil
}

func (t *{{.CamlName}}) Create(ctx context.Context, info *{{.UpperCamlName}}, opts ...model.Option) (err error) {
	build, f, err := model.New().Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(info)}, opts...)...)
	if err != nil {
		return err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	_, err = build.Insert()
	if err != nil {
		return err
	}
	return nil
}

func (t *{{.CamlName}}) Update(ctx context.Context, con model.IQuery, info *{{.UpperCamlName}}, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(info)}, opts...)...)
	if err != nil {
		return err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	result, err := build.Update()
	if err != nil {
		return err
	}
	if build.Result != nil {
		build.Gain(result.RowsAffected(), result.RowsReturned())
	}
	return nil
}

func (t *{{.CamlName}}) Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&{{.UpperCamlName}}{})}, opts...)...)
	if err != nil {
		return err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	_, err = build.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (t *{{.CamlName}}) Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *{{.UpperCamlName}}, err error) {
	detail := new({{.UpperCamlName}})
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(detail)}, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	if err = build.First(); err != nil {
		return nil, err
	}
	return detail, nil
}

func (t *{{.CamlName}}) List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*{{.UpperCamlName}}, err error) {
	list := make([]*{{.UpperCamlName}}, 0)
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&list)}, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	if err = build.Select(); err != nil {
		return nil, err
	}
	return list, nil
}

func (t *{{.CamlName}}) Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&{{.UpperCamlName}}{})}, opts...)...)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	total, err := build.Count()
	if err != nil {
		return 0, err
	}
	return int32(total), nil
}

`
var ConfTemplate2 = `// package {{.Package}}
// nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import (
	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
)

// ID 条件
func ID(id int32) model.QueryOption {
	return func(query *pg.Query) {
		query.Where("id = ?", id)
	}
}

// GetCommonQuery 获取公共查询条件
func GetCommonQuery(params http.BaseRequest) model.QueryOption {
	prefix := model.GetTableName((*{{.UpperCamlName}})(nil))
	return func(query *pg.Query) {
		model.GetCommonQuery(prefix, params)(query)
		for _, search := range params.Search {
			switch search.Key {
			default:
				search.Key = fmt.Sprintf("%s.%s", prefix, search.Key)
				// 通用字段处理
				model.DealSearch(query, search)
			}
		}
	}
}

// WithColumns 获取那些字段
func WithColumns(columns ...string) model.QueryOption {
	return func(query *pg.Query) {
		query.Column(columns...)
	}
}

`

var ConfTemplate3 = `// package {{.Package}} 
// nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import (
	"github.com/thoohv5/person/internal/model"
)

func init() {
	model.Register((*{{.UpperCamlName}})(nil))
}

// {{.UpperCamlName}} 模块
type {{.UpperCamlName}} struct {
	model.BaseModel
	tableName struct{} {{.Backquote}}pg:"{{.Name}}s,discard_unknown_columns"{{.Backquote}}

}
`
