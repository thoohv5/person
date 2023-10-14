// Package service 服务
package service

import (
	"context"
	"errors"

	"github.com/thoohv5/person/internal/provide/http"
	"github.com/thoohv5/person/pkg/util/errgroup"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/app/interface/demo-interface/api/http/request"
	"github.com/thoohv5/person/app/interface/demo-interface/api/http/response"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/repository"
	rDemo "github.com/thoohv5/person/app/interface/demo-interface/internal/repository/demo"
	"github.com/thoohv5/person/internal/code"
	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/internal/provide/nats"
	"github.com/thoohv5/person/pkg/log"
)

// demo 模板
type demo struct {
	log      log.Logger
	conf     config.Config
	data     repository.IRepository
	producer nats.IProducer
}

// NewDemo 创建模板
func NewDemo(
	logger log.Logger,
	conf config.Config,
	data repository.IRepository,
) IDemo {
	return &demo{
		log:  logger,
		conf: conf,
		data: data,
	}
}

// IDemo 模板
type IDemo interface {
	// Create 创建
	Create(ctx context.Context, param *request.DemoCreate) error
	// Update 更新
	Update(ctx context.Context, uriParam *request.DemoPkID, param *request.DemoUpdate) error
	// List 列表
	List(ctx context.Context, param *request.DemoList) ([]*response.DemoEntity, int32, error)
	// Detail 详情
	Detail(ctx context.Context, param *request.DemoPkID) (*rDemo.Demo, error)
	// Delete 删除
	Delete(ctx context.Context, param *request.DemoPkID) error
}

// Create 创建
func (d *demo) Create(ctx context.Context, param *request.DemoCreate) error {
	// 创建
	if err := d.data.GetDemo().Create(ctx, &rDemo.Demo{
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
func (d *demo) Update(ctx context.Context, uriParam *request.DemoPkID, param *request.DemoUpdate) error {
	// 更新
	if err := d.data.GetDemo().Update(ctx, model.Where(
		rDemo.ID(uriParam.ID),
	), &rDemo.Demo{
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
func (d *demo) List(ctx context.Context, param *request.DemoList) ([]*response.DemoEntity, int32, error) {
	ret := make([]*response.DemoEntity, 0)
	total := int32(0)

	// 条件
	condition := model.Where(
		rDemo.GetCommonQuery(http.BaseRequest{
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
		cnt, err := d.data.GetDemo().Count(ctx, condition)
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
		list, err := d.data.GetDemo().List(ctx, condition)
		if err != nil {
			d.log.Errorc(ctx, "[demo][List] list db error", logger.FieldError(err))
			return code.ErrOpDB
		}
		// model转化entity
		entityList := make([]*response.DemoEntity, 0, len(list))
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
func (d *demo) Detail(ctx context.Context, param *request.DemoPkID) (*rDemo.Demo, error) {
	// 详情
	data, err := d.data.GetDemo().Detail(ctx, model.Where(
		rDemo.ID(param.ID),
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
func (d *demo) Delete(ctx context.Context, param *request.DemoPkID) error {
	// 删除
	if err := d.data.GetDemo().Delete(ctx, model.Where(
		rDemo.ID(param.ID),
	)); err != nil {
		d.log.Errorc(ctx, "[demo Delete] delete db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// 转换
func (d *demo) toEntity(item *rDemo.Demo) *response.DemoEntity {
	return &response.DemoEntity{
		ID:         item.ID,
		Name:       item.Name,
		State:      item.State,
		Age:        item.Age,
		CreateTime: item.CreatedTime,
	}
}
