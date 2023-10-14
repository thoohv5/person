package template

var Service = `package service

import (
	"context"
	"fmt"

	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/api/http/request"
	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/api/http/response"
	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/internal/repository"
	r{{.UpperCamlName}} "github.com/thoohv5/person/app/interface/{{.ProjectName}}/internal/repository/{{.Package}}"
	"github.com/thoohv5/person/internal/code"
	"github.com/thoohv5/person/internal/http"
	"github.com/thoohv5/person/internal/model"
	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

// {{.CamlName}} struct
type {{.CamlName}} struct {
	log  log.Logger
	data repository.IRepository
}

// New{{.UpperCamlName}} 创建
func New{{.UpperCamlName}}(logger log.Logger, data repository.IRepository) I{{.UpperCamlName}} {
	return &{{.CamlName}}{
		log:  logger,
		data: data,
	}
}

// I{{.UpperCamlName}} 接口
type I{{.UpperCamlName}} interface {
	All(ctx context.Context, params *request.{{.UpperCamlName}}All) (*http.ListResponse, error)
	List(ctx context.Context, params *request.{{.UpperCamlName}}List) (*http.PageResponse, error)
	Create(ctx context.Context, params *request.{{.UpperCamlName}}Create) error
	Update(ctx context.Context, uriParams *request.{{.UpperCamlName}}PkID, params *request.{{.UpperCamlName}}Update) error
	Detail(ctx context.Context, params *request.{{.UpperCamlName}}PkID) (*response.{{.UpperCamlName}}Entity, error)
	Delete(ctx context.Context, params *request.{{.UpperCamlName}}Delete) error
}

// All 获取全部数据
func (d *{{.CamlName}}) All(ctx context.Context, params *request.{{.UpperCamlName}}All) (*http.ListResponse, error) {
	list, err := d.data.Get{{.UpperCamlName}}().List(ctx, model.Where(
		model.Common(model.BaseRequest{
			Sorts: []string{"age desc", "name"},
		}),
	))
	if err != nil {
		d.log.Errorc(ctx, "[LIST-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	// model转化entity
	entityList := make([]*response.{{.UpperCamlName}}Entity, 0)
	for _, item := range list {
		entityList = append(entityList, d.toEntity(ctx, item))
	}
	return &http.ListResponse{List: entityList}, err
}

// List 获取分页数据
func (d *{{.CamlName}}) List(ctx context.Context, params *request.{{.UpperCamlName}}List) (*http.PageResponse, error) {
	result := &http.PageResponse{
		List:  make([]*response.{{.UpperCamlName}}Entity, 0),
		PageNum: params.PageNum,
		PageSize: params.PageSize,
		Total: 0,
	}
	condition := model.Where(
		r{{.UpperCamlName}}.GetCommonQuery(params.BaseRequest),
	)

	total, err := d.data.Get{{.UpperCamlName}}().Count(ctx, model.New())
	if err != nil {
		d.log.Errorc(ctx, "[COUNT-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	result.Total = total
	start := (params.PageNum - 1) * params.PageSize
	if total <= start {
		return result, nil
	}

	list, err := d.data.Get{{.UpperCamlName}}().List(ctx, condition)
	if err != nil {
		d.log.Errorc(ctx, "[LIST-ERROR] db error", logger.FieldError(err))
		return nil, code.ErrOpDB
	}
	// model转化entity
	entityList := make([]*response.{{.UpperCamlName}}Entity, 0, len(list))
	for _, item := range list {
		entityList = append(entityList, d.toEntity(ctx, item))
	}
	result.List = entityList
	return result, nil
}

// Create 创建
func (d *{{.CamlName}}) Create(ctx context.Context, params *request.{{.UpperCamlName}}Create) error {
	_, err := d.data.Get{{.UpperCamlName}}().Detail(ctx, model.Where(
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return code.ErrOpDB
		}
	} else {
		return code.ErrDataExist
	}

	data := &r{{.UpperCamlName}}.{{.UpperCamlName}}{}
	err = d.data.Get{{.UpperCamlName}}().Create(ctx, data)
	if err != nil {
		d.log.Infoc(ctx, "[CREATE-ERROR] db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// Update 更新
func (d *{{.CamlName}}) Update(ctx context.Context, uriParams *request.{{.UpperCamlName}}PkID, params *request.{{.UpperCamlName}}Update) error {
	_, err := d.data.Get{{.UpperCamlName}}().Detail(ctx, model.Where(
		r{{.UpperCamlName}}.ID(uriParams.ID),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return code.ErrOpDB
		}
		return code.ErrDataNotExist
	}

	data := &r{{.UpperCamlName}}.{{.UpperCamlName}}{}
	err = d.data.Get{{.UpperCamlName}}().Update(ctx, model.Where(
		r{{.UpperCamlName}}.ID(uriParams.ID),
	), data)
	if err != nil {
		d.log.Infoc(ctx, "[UPDATE-ERROR] db error", logger.FieldError(err))
		return code.ErrOpDB
	}

	return nil
}

// Detail 获取详情
func (d *{{.CamlName}}) Detail(ctx context.Context, params *request.{{.UpperCamlName}}PkID) (*response.{{.UpperCamlName}}Entity, error) {
	data, err := d.data.Get{{.UpperCamlName}}().Detail(ctx, model.Where(
		r{{.UpperCamlName}}.ID(params.ID),
	))

	if err != nil {
		if model.IsNotErrNoRows(err) {
			d.log.Errorc(ctx, "[DETAIL-ERROR] db error", logger.FieldError(err))
			return nil, code.ErrOpDB
		}
		return nil, code.ErrDataNotExist
	}
	return d.toEntity(ctx, data), nil
}

// Delete 删除
func (d *{{.CamlName}}) Delete(ctx context.Context, params *request.{{.UpperCamlName}}Delete) error {
	// 自己实现
	return nil
}

// toEntity 自定义实体转换
func (d *{{.CamlName}}) toEntity(ctx context.Context, model *r{{.UpperCamlName}}.{{.UpperCamlName}}) *response.{{.UpperCamlName}}Entity {
	return &response.{{.UpperCamlName}}Entity{}
}
`
