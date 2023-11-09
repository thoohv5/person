package {{ .Name }}

import (
	"context"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
)

type {{ .Name }} struct {
	db pg.DBI
}

// I{{ U .Name }} 标准
type I{{ U .Name }} interface {
	Create(ctx context.Context, info *{{ U .Name }}, opts ...model.Option) (err error)
	Update(ctx context.Context, con model.IQuery, info *{{ U .Name }}, opts ...model.Option) (err error)
	Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error)
	Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *{{ U .Name }}, err error)
	List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*{{ U .Name }}, err error)
	Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error)
}

// New 创建
func New(db pg.DBI) (I{{ U .Name }}, error) {
	return &{{ .Name }}{
		db: db,
	}, nil
}

func (t *{{ .Name }}) Create(ctx context.Context, info *{{ U .Name }}, opts ...model.Option) (err error) {
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

func (t *{{ .Name }}) Update(ctx context.Context, con model.IQuery, info *{{ U .Name }}, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(info)}, opts...)...)
	if err != nil {
		return err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	_, err = build.Update()
	if err != nil {
		return err
	}
	return nil
}

func (t *{{ .Name }}) Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&{{ U .Name }}{})}, opts...)...)
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

func (t *{{ .Name }}) Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *{{ U .Name }}, err error) {
	detail := new({{ U .Name }})
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

func (t *{{ .Name }}) List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*{{ U .Name }}, err error) {
	list := make([]*{{ U .Name }}, 0)
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

func (t *{{ .Name }}) Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&{{ U .Name }}{})}, opts...)...)
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
