package {{ .Name }}

import (
	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/internal/provide/http"

	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/api/config"
	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/api/http/request"
	"github.com/thoohv5/person/app/interface/{{ .ProjectName }}/internal/service"
	"github.com/thoohv5/person/pkg/log"
)

// {{ U .Name }} {{ .Remark }}
type {{ U .Name }} struct {
	conf    config.Config
	logger  log.Logger
	service service.I{{ U .Name }}
}

// New{{ U .Name }} 创建
func New{{ U .Name }}(
	conf config.Config,
	logger log.Logger,
	service service.I{{ U .Name }},
) *{{ U .Name }} {
	return &{{ U .Name }}{
		conf:    conf,
		logger:  logger,
		service: service,
	}
}

// Create 创建
//
//	@Summary		创建{{ .Remark }}
//	@Description	描述
//	@Tags			{{ .Remark }}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			req	body		request.{{ U .Name }}Create	true	"请求参数"
//	@Success		200	{object}	http.Response{}		"返回值"
//	@Router			/{{ .ProjectName }}/{{ .Name }} [post]
func (d *{{ U .Name }}) Create(gtx *gin.Context) {
	req := &request.{{ U .Name }}Create{}

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
//
//	@Summary		修改{{ .Remark }}
//	@Description	描述
//	@Tags			{{ .Remark }}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int					true	"{{ .Remark }}ID"
//	@Param			req	body		request.{{ U .Name }}Update	true	"请求参数"
//	@Success		200	{object}	http.Response{}		"返回值"
//	@Router			/{{ .ProjectName }}/{{ .Name }}/{id} [put]
func (d *{{ U .Name }}) Update(gtx *gin.Context) {
	uriReq := &request.{{ U .Name }}PkID{}
	if err := gtx.BindUri(uriReq); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	req := &request.{{ U .Name }}Update{}
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

// List 获取分页列表
//
//	@Summary		获取{{ .Remark }}分页数据
//	@Description	描述
//	@Tags			{{ .Remark }}
//	@Security		ApiKeyAuth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			req	query		request.{{ U .Name }}List													true	"请求参数"
//	@Success		200	{object}	http.Response{data=http.PageResponse{list=[]response.{{ U .Name }}Entity}}	"返回值"
//	@Router			/{{ .ProjectName }}/{{ .Name }} [get]
func (d *{{ U .Name }}) List(gtx *gin.Context) {
	req := &request.{{ U .Name }}List{}

	if err := gtx.Bind(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	list, total, err := d.service.List(gtx, req)
	if err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, &http.PageResponse{
		List:  list,
		Total: total,
	})
}

// Detail 获取详细数据
//
//	@Summary		获取{{ .Remark }}详细数据
//	@Description	描述
//	@Tags			{{ .Remark }}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int										true	"{{ .Remark }}ID"
//	@Success		200	{object}	http.Response{data=response.{{ U .Name }}Entity}	"返回值"
//	@Router			/{{ .ProjectName }}/{{ .Name }}/{id} [get]
func (d *{{ U .Name }}) Detail(gtx *gin.Context) {
	req := &request.{{ U .Name }}PkID{}
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
//
//	@Summary		删除{{ .Remark }}
//	@Description	描述
//	@Tags			{{ .Remark }}
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int				true	"{{ .Remark }}ID"
//	@Success		200	{object}	http.Response{}	"返回值"
//	@Router			/{{ .ProjectName }}/{{ .Name }}/{id} [delete]
func (d *{{ U .Name }}) Delete(gtx *gin.Context) {
	req := &request.{{ U .Name }}PkID{}
	if err := gtx.BindUri(req); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	if err := d.service.Delete(gtx, req); err != nil {
		http.Fail(gtx, err)
		return
	}

	http.Success(gtx, nil)
}
