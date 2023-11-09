// Package service 服务
package service

import (
	"context"
	"errors"

	"github.com/thoohv5/person/internal/provide/http"
	"github.com/thoohv5/person/pkg/util/errgroup"

	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/api/config"
	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/api/http/request"
	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/api/http/response"
	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/internal/repository"
	r{{ U .Name }} "github.com/thoohv5/person/app/interface/{{ .ProjectName }}/internal/repository/{{ .Name }}"
	"github.com/thoohv5/person/internal/code"
	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

// {{ .Name }} {{ .Remark }}
type {{ .Name }} struct {
	log      log.Logger
	conf     config.Config
	data     repository.IRepository
}

// {{ U .Name }}Demo 创建{{ .Remark }}
func New{{ U .Name }}(
	logger log.Logger,
	conf config.Config,
	data repository.IRepository,
) I{{ U .Name }} {
	return &{{ .Name }}{
		log:  logger,
		conf: conf,
		data: data,
	}
}

// I{{ U .Name }} 模板
type I{{ U .Name }} interface {
	// Create 创建
	Create(ctx context.Context, param *request.{{ U .Name }}Create) error
	// Update 更新
	Update(ctx context.Context, uriParam *request.{{ U .Name }}PkID, param *request.{{ U .Name }}Update) error
	// List 列表
	List(ctx context.Context, param *request.{{ U .Name }}List) ([]*response.{{ U .Name }}Entity, int32, error)
	// Detail 详情
	Detail(ctx context.Context, param *request.{{ U .Name }}PkID) (*r{{ U .Name }}.{{ U .Name }}, error)
	// Delete 删除
	Delete(ctx context.Context, param *request.{{ U .Name }}PkID) error
}

// Create 创建
func (d *{{ .Name }}) Create(ctx context.Context, param *request.{{ U .Name }}Create) error {
	// 创建
	if err := d.data.Get{{ U .Name }}().Create(ctx, &r{{ U .Name }}.{{ U .Name }}{
		Name:  param.Name,
		Age:   param.Age,
		State: param.State,
	}); err != nil {
		d.log.Errorc(ctx, "[demo Create] create db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// Update 更新
func (d *{{ .Name }}) Update(ctx context.Context, uriParam *request.{{ U .Name }}PkID, param *request.{{ U .Name }}Update) error {
	// 更新
	if err := d.data.Get{{ U .Name }}().Update(ctx, model.Where(
		r{{ U .Name }}.ID(uriParam.ID),
	), &r{{ U .Name }}.{{ U .Name }}{
		Name:  param.Name,
		Age:   param.Age,
		State: param.State,
	}); err != nil {
		d.log.Errorc(ctx, "[demo Update] update db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// List 获取分页数据
func (d *{{ .Name }}) List(ctx context.Context, param *request.{{ U .Name }}List) ([]*response.{{ U .Name }}Entity, int32, error) {
	ret := make([]*response.{{ U .Name }}Entity, 0)
	total := int32(0)

	// 条件
	condition := model.Where(
		r{{ U .Name }}.GetCommonQuery(http.BaseRequest{
			PageNum:   param.PageNum,
			PageSize:  param.PageSize,
			SortField: param.SortField,
			SortOrder: param.SortOrder,
			Search:    param.Search,
		}),
		// todo 条件
	)

	// 并发
	eg := errgroup.WithContext(ctx)

	// 总数
	eg.Go(func(ctx context.Context) error {
		if param.ExclusiveTotal {
			return nil
		}
		cnt, err := d.data.Get{{ U .Name }}().Count(ctx, condition)
		if err != nil {
			d.log.Errorc(ctx, "[demo][List] count db error", logger.FieldError(err))
			return code.ErrOpDB
		}
		if cnt == 0 {
			return code.ErrOpCancel
		}
		total = cnt
		return nil
	})

	// 列表
	eg.Go(func(ctx context.Context) error {
		if param.ExclusiveList {
			return nil
		}
		list, err := d.data.Get{{ U .Name }}().List(ctx, condition)
		if err != nil {
			d.log.Errorc(ctx, "[demo][List] list db error", logger.FieldError(err))
			return code.ErrOpDB
		}
		// model转化entity
		entityList := make([]*response.{{ U .Name }}Entity, 0, len(list))
		for _, item := range list {
			entityList = append(entityList, d.toEntity(item))
		}
		ret = entityList
		return nil
	})

	// 等待
	if err := eg.Wait(); err != nil && !errors.Is(err, code.ErrOpCancel) {
		return nil, 0, err
	}

	return ret, total, nil
}

// Detail 创建
func (d *{{ .Name }}) Detail(ctx context.Context, param *request.{{ U .Name }}PkID) (*r{{ U .Name }}.{{ U .Name }}, error) {
	// 详情
	data, err := d.data.Get{{ U .Name }}().Detail(ctx, model.Where(
		r{{ U .Name }}.ID(param.ID),
	))
	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[demo Detail] detail db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
		return nil, code.ErrDataNotExist
	}
	return data, nil
}

// Delete 删除
func (d *{{ .Name }}) Delete(ctx context.Context, param *request.{{ U .Name }}PkID) error {
	// 删除
	if err := d.data.Get{{ U .Name }}().Delete(ctx, model.Where(
		r{{ U .Name }}.ID(param.ID),
	)); err != nil {
		d.log.Errorc(ctx, "[demo Delete] delete db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// 转换
func (d *{{ .Name }}) toEntity(item *r{{ U .Name }}.{{ U .Name }}) *response.{{ U .Name }}Entity {
	return &response.{{ U .Name }}Entity{
		ID:         item.ID,
		Name:       item.Name,
		State:      item.State,
		Age:        item.Age,
		CreateTime: item.CreatedTime,
	}
}
