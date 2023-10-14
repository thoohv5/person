package demo

import (
	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/internal/provide/http"

	"github.com/thoohv5/person/app/interface/demo-interface/api/config"
	"github.com/thoohv5/person/app/interface/demo-interface/api/http/request"
	"github.com/thoohv5/person/app/interface/demo-interface/internal/service"
	"github.com/thoohv5/person/pkg/log"
)

// Demo 模板
type Demo struct {
	conf    config.Config
	logger  log.Logger
	service service.IDemo
}

// NewDemo 创建
func NewDemo(
	conf config.Config,
	logger log.Logger,
	service service.IDemo,
) *Demo {
	return &Demo{
		conf:    conf,
		logger:  logger,
		service: service,
	}
}

// Create 创建
//
//	@Summary		创建模板
//	@Description	描述
//	@Tags			模板
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			req	body		request.DemoCreate	true	"请求参数"
//	@Success		200	{object}	http.Response{}		"返回值"
//	@Router			/demo-interface/demo [post]
func (d *Demo) Create(gtx *gin.Context) {
	req := &request.DemoCreate{}

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
//	@Summary		修改模板
//	@Description	描述
//	@Tags			模板
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int					true	"demo ID"
//	@Param			req	body		request.DemoUpdate	true	"请求参数"
//	@Success		200	{object}	http.Response{}		"返回值"
//	@Router			/demo-interface/demo/{id} [put]
func (d *Demo) Update(gtx *gin.Context) {
	uriReq := &request.DemoPkID{}
	if err := gtx.BindUri(uriReq); err != nil {
		http.BadRequest(gtx, err)
		return
	}

	req := &request.DemoUpdate{}
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
//	@Summary		获取模板分页数据
//	@Description	描述
//	@Tags			模板
//	@Security		ApiKeyAuth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			req	query		request.DemoList													true	"请求参数"
//	@Success		200	{object}	http.Response{data=http.PageResponse{list=[]response.DemoEntity}}	"返回值"
//	@Router			/demo-interface/demo [get]
func (d *Demo) List(gtx *gin.Context) {
	req := &request.DemoList{}

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
//	@Summary		获取模板详细数据
//	@Description	描述
//	@Tags			模板
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int										true	"demo ID"
//	@Success		200	{object}	http.Response{data=response.DemoEntity}	"返回值"
//	@Router			/demo-interface/demo/{id} [get]
func (d *Demo) Detail(gtx *gin.Context) {
	req := &request.DemoPkID{}
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
//	@Summary		删除模板
//	@Description	描述
//	@Tags			模板
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int				true	"demo ID"
//	@Success		200	{object}	http.Response{}	"返回值"
//	@Router			/demo-interface/demo/{id} [delete]
func (d *Demo) Delete(gtx *gin.Context) {
	req := &request.DemoPkID{}
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
