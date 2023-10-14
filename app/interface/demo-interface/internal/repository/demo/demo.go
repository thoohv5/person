package demo

import (
	"context"

	"github.com/go-pg/pg/v10"

	"github.com/thoohv5/person/internal/model"
)

type demo struct {
	db pg.DBI
}

// IDemo 标准
type IDemo interface {
	Create(ctx context.Context, info *Demo, opts ...model.Option) (err error)
	Update(ctx context.Context, con model.IQuery, info *Demo, opts ...model.Option) (err error)
	Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error)
	Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *Demo, err error)
	List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*Demo, err error)
	Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error)
}

// New 创建
func New(db pg.DBI) (IDemo, error) {
	return &demo{
		db: db,
	}, nil
}

func (t *demo) Create(ctx context.Context, info *Demo, opts ...model.Option) (err error) {
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

func (t *demo) Update(ctx context.Context, con model.IQuery, info *Demo, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(info)}, opts...)...)
	if err != nil {
		return err
	}
	defer func() {
		if err = (*f)(err); err != nil {
			return
		}
	}()
	_, err = build.UpdateNotZero()
	if err != nil {
		return err
	}
	return nil
}

func (t *demo) Delete(ctx context.Context, con model.IQuery, opts ...model.Option) (err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&Demo{})}, opts...)...)
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

func (t *demo) Detail(ctx context.Context, con model.IQuery, opts ...model.Option) (ret *Demo, err error) {
	detail := new(Demo)
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

func (t *demo) List(ctx context.Context, con model.IQuery, opts ...model.Option) (ret []*Demo, err error) {
	list := make([]*Demo, 0)
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

func (t *demo) Count(ctx context.Context, con model.IQuery, opts ...model.Option) (cnt int32, err error) {
	build, f, err := con.Build(ctx, append([]model.Option{model.WithDb(t.db), model.WithModel(&Demo{})}, opts...)...)
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
