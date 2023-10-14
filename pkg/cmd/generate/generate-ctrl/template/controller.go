package template

var Controller = `package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/api/http/request"
	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/api/conf"
	"github.com/thoohv5/person/app/interface/{{.ProjectName}}/internal/service"
	"github.com/thoohv5/person/internal/http"
	"github.com/thoohv5/person/pkg/log"
)

// {{.UpperCamlName}} 名称
type {{.UpperCamlName}} struct {
	conf    conf.Config
	logger  log.Logger
	service service.I{{.UpperCamlName}}
}

// New{{.UpperCamlName}} 创建
func New{{.UpperCamlName}}(
	conf conf.Config,
	logger log.Logger,
	service service.I{{.UpperCamlName}},
) *{{.UpperCamlName}} {
	return &{{.UpperCamlName}}{
		conf:    conf,
		logger:  logger,
		service: service,
	}
}

// Create 创建
//	@Summary		创建{{.NameTrans}}
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			req	body		request.{{.UpperCamlName}}Create						true	"请求参数"
//	@Router			/{{.ProjectName}}/{{.MidlineName}} [post]
func (d *{{.UpperCamlName}}) Create(gtx *gin.Context) {
	req := &request.{{.UpperCamlName}}Create{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	if err := d.service.Create(gtx, req); err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, nil)
}

// Update 修改
//	@Summary		修改{{.NameTrans}}
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int										true	"{{.UpperCamlName}} ID"
//	@Param			req	body		request.{{.UpperCamlName}}Update						true	"请求参数"
//	@Router			/{{.ProjectName}}/{{.MidlineName}}/{id} [put]
func (d *{{.UpperCamlName}}) Update(gtx *gin.Context) {
	uriReq := &request.{{.UpperCamlName}}PkID{}
	if err := gtx.BindUri(uriReq); err != nil {
		http.BadRequest(gtx, err)
		return
	}
	req := &request.{{.UpperCamlName}}Update{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	if err := d.service.Update(gtx, uriReq, req); err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, nil)
}

// All 获取列表
//	@Summary		获取{{.NameTrans}}
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			req	query		request.{{.UpperCamlName}}All														true	"请求参数"
//	@Success		200	{object}	http.Response{data=http.ListResponse{list=[]response.{{.UpperCamlName}}Entity}}	"返回值"
//	@Router			/{{.ProjectName}}/{{.MidlineName}}/all [get]
func (d *{{.UpperCamlName}}) All(gtx *gin.Context) {
	req := &request.{{.UpperCamlName}}All{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	resp, err := d.service.All(gtx, req)
	if err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, resp)
}

// List 获取分页列表
//	@Summary		获取{{.NameTrans}}分页数据
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			req	query		request.{{.UpperCamlName}}List													true	"请求参数"
//	@Success		200	{object}	http.Response{data=http.ListResponse{list=[]response.{{.UpperCamlName}}Entity}}	"返回值"
//	@Router			/{{.ProjectName}}/{{.MidlineName}} [get]
func (d *{{.UpperCamlName}}) List(gtx *gin.Context) {
	req := &request.{{.UpperCamlName}}List{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	resp, err := d.service.List(gtx, req)
	if err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, resp)
}

// Detail 获取详细数据
//	@Summary		获取{{.NameTrans}}详细数据
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int										true	"{{.UpperCamlName}} ID"
//	@Success		200	{object}	http.Response{data=response.{{.UpperCamlName}}Entity}	"返回值"
//	@Router			/{{.ProjectName}}/{{.MidlineName}}/{id} [get]
func (d *{{.UpperCamlName}}) Detail(gtx *gin.Context) {
	req := &request.{{.UpperCamlName}}PkID{}
	if err := gtx.BindUri(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}
	resp, err := d.service.Detail(gtx, req)
	if err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, resp)
}

// Delete 删除
//	@Summary		删除{{.NameTrans}}
//	@Description	描述
//	@Tags			{{.NameTrans}}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			req	body		request.{{.UpperCamlName}}Delete	true	"请求参数"
//	@Router			/{{.ProjectName}}/{{.MidlineName}} [delete]
func (d *{{.UpperCamlName}}) Delete(gtx *gin.Context) {
	req := &request.{{.UpperCamlName}}Delete{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	if err := d.service.Delete(gtx, req); err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, nil)
}
`
