// Package service 服务
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/thoohv5/person/internal/provide/http"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/app/interface/demo-interface/api/http/request"
	"github.com/thoohv5/person/app/interface/demo-interface/api/http/response"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/repository"
	rdm "github.com/thoohv5/person/app/interface/demo-interface/internal/repository/demo"
	"github.com/thoohv5/person/internal/code"
	"github.com/thoohv5/person/internal/localize"
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
	All(ctx context.Context, params *request.DemoAll) (*http.ListResponse, error)
	List(ctx context.Context, params *request.DemoList) (*http.PageResponse, error)
	Create(ctx context.Context, params *request.DemoCreate) (*rdm.Demo, error)
	Update(ctx context.Context, uriParams *request.DemoPkID, params *request.DemoUpdate) (*rdm.Demo, error)
	Detail(ctx context.Context, params *request.DemoPkID) (*rdm.Demo, error)
	Delete(ctx context.Context, params *request.DemoPkID) (*rdm.Demo, error)
}

// All 获取全部数据
func (d *demo) All(ctx context.Context, params *request.DemoAll) (*http.ListResponse, error) {
	d.log.Infoc(ctx, localize.Translate(ctx, "欢迎您!"))

	list, err := d.data.GetDemo().List(ctx, model.Where(
		model.Common(model.BaseRequest{
			Sorts: []string{"age desc", "name"},
		}),
		rdm.State(params.State),
	))
	if err != nil {
		d.log.Errorc(ctx, "[LIST-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	// model转化entity
	entityList := make([]*response.DemoEntity, 0)
	for _, item := range list {
		entityList = append(entityList, &response.DemoEntity{
			ID:         item.ID,
			Name:       item.Name,
			State:      item.State,
			Age:        item.Age,
			CreateTime: item.CreatedTime,
		})
	}
	return &http.ListResponse{List: entityList}, err
}

// List 获取分页数据
func (d *demo) List(ctx context.Context, params *request.DemoList) (*http.PageResponse, error) {
	d.log.Errorc(ctx, fmt.Sprintf("这里测试一下普通业务日志:[%s]", "test"), logger.FieldString("test", "测试"))
	result := &http.PageResponse{
		List:  make([]*response.DemoEntity, 0),
		Total: 0,
	}
	limit := params.PageSize
	start := (params.PageNum - 1) * params.PageSize
	condition := model.Where(
		model.Common(model.BaseRequest{
			Sorts: []string{fmt.Sprintf("%s %s", params.SortField, params.SortOrder)},
			Limit: limit,
			Start: start,
		}),
	)

	total, err := d.data.GetDemo().Count(ctx, model.New())
	if err != nil {
		d.log.Errorc(ctx, "[COUNT-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	result.Total = total
	if total <= start {
		return result, nil
	}

	list, err := d.data.GetDemo().List(ctx, condition)
	if err != nil {
		d.log.Errorc(ctx, "[LIST-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	// model转化entity
	entityList := make([]*response.DemoEntity, 0, len(list))
	for _, item := range list {
		entityList = append(entityList, &response.DemoEntity{
			ID:         item.ID,
			Name:       item.Name,
			State:      item.State,
			Age:        item.Age,
			CreateTime: item.CreatedTime,
		})
	}
	result.List = entityList
	return result, nil
}

// Create 创建
func (d *demo) Create(ctx context.Context, params *request.DemoCreate) (*rdm.Demo, error) {
	_, err := d.data.GetDemo().Detail(ctx, model.Where(
		rdm.Name(params.Name),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
	} else {
		return nil, code.ErrDataExist
	}

	data := &rdm.Demo{}
	data.Name = params.Name
	data.State = params.State
	data.Age = params.Age
	err = d.data.GetDemo().Create(ctx, data)
	if err != nil {
		d.log.Infoc(ctx, "[CREATE-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}

	if err := d.data.RunInTransaction(ctx, func(ctx context.Context, re repository.IRepository, op ...model.Option) error {
		if err := re.GetDemo().Create(ctx, &rdm.Demo{
			Name: fmt.Sprintf("xx_%d", time.Now().Unix()),
		}, op...); err != nil {
			return err
		}

		if err := re.GetDemo().Update(ctx, model.Where(rdm.ID(1)), &rdm.Demo{
			State: 1,
		}, op...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return data, nil
}

// Update 更新
func (d *demo) Update(ctx context.Context, uriParams *request.DemoPkID, params *request.DemoUpdate) (*rdm.Demo, error) {
	_, err := d.data.GetDemo().Detail(ctx, model.Where(
		rdm.ID(uriParams.ID),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
		return nil, code.ErrDataNotExist
	}

	data := &rdm.Demo{}
	data.Name = params.Name
	data.State = params.State
	data.Age = params.Age
	err = d.data.GetDemo().Update(ctx, model.Where(
		rdm.ID(uriParams.ID),
	), data)
	if err != nil {
		d.log.Infoc(ctx, "[UPDATE-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}

	return data, nil
}

// Detail 创建
func (d *demo) Detail(ctx context.Context, params *request.DemoPkID) (*rdm.Demo, error) {
	data, err := d.data.GetDemo().Detail(ctx, model.Where(
		rdm.ID(params.ID),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
		return nil, code.ErrDataNotExist
	}
	return data, nil
}

// Delete 删除
func (d *demo) Delete(ctx context.Context, params *request.DemoPkID) (*rdm.Demo, error) {
	data, err := d.data.GetDemo().Detail(ctx, model.Where(
		rdm.ID(params.ID),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
		return nil, code.ErrDataNotExist
	}

	err = d.data.GetDemo().Delete(ctx, model.Where(
		rdm.ID(params.ID),
	))
	if err != nil {
		d.log.Errorc(ctx, "[DELETE-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}

	return data, nil
}
