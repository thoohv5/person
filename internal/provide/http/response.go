package http

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/thoohv5/person/internal/localize"

	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/internal/code"
	"github.com/thoohv5/person/internal/validate"
)

// NewlineTag 换行符
const NewlineTag = "<\\br>"

// Response 返回结构
type Response struct {
	// 错误码
	Code string `json:"code"`
	// 提示信息
	Msg string `json:"msg"`
	// 提示详情
	Detail interface{} `json:"detail"`
	// 数据
	Data interface{} `json:"data"`
}

// ListResponse 列表返回结构
type ListResponse struct {
	// 数据
	List interface{} `json:"list"`
}

// PageResponse 列表返回结构
type PageResponse struct {
	// 数据
	List interface{} `json:"list"`
	// 页码
	PageNum int32 `json:"pageNum"`
	// 每页数量
	PageSize int32 `json:"pageSize"`
	// 总数
	Total int32 `json:"total"`
}

// Success 成功
func Success(ctx *gin.Context, data interface{}) {
	c := code.Success
	base(ctx, c.Status(), c.Code(), c.Error(), data)
}

// Fail 失败
func Fail(ctx *gin.Context, err error) {
	c := &code.Code{}
	if errors.As(err, &c) {
		base(ctx, c.Status(), c.Code(), c.Error(), c.Detail())
		return
	}
	var e error
	if errors.As(err, &e) {
		c := code.ErrCommon
		base(ctx, c.Status(), c.Code(), e.Error(), nil)
	}
}

// BadRequest 参数异常
func BadRequest(ctx *gin.Context, err error) {
	c := code.ErrParamInvalid
	msg := validate.Translate(err)
	switch mmsg := msg.(type) {
	case map[string]string:
		str := ""
		for _, s2 := range mmsg {
			str = fmt.Sprintf("%s,%s", str, s2)
		}
		base(ctx, c.Status(), c.Code(), strings.TrimLeft(str, ","), nil)
	case string:
		base(ctx, c.Status(), c.Code(), mmsg, nil)
	}
}

// base 基础函数
func base(ctx *gin.Context, httpCode, msgCode int, msg string, data interface{}) {
	var detail interface{}
	if code.Success.Code() != msgCode {
		// 处理批量删除时，数据不存在场景
		switch d := data.(type) {
		case []string:
			str := ""
			for _, s := range d {
				str = fmt.Sprintf("%s%s%s", str, NewlineTag, s)
			}
			msg = fmt.Sprintf("%s%s", getMsg(ctx, msgCode), str)
		case map[string]string:
			// 创建网络时，数据存在冲突
			str := ""
			for ip, net := range d["standard"] {
				str = fmt.Sprintf("%s%s%s", str, NewlineTag, localize.Translate(ctx, "IP[%s]被标准网络[%s]包含", ip, net))
			}
			msg = fmt.Sprintf("%s%s", getMsg(ctx, msgCode), str)
		}
		detail = data
		data = gin.H{}
	}
	if data == nil || fmt.Sprintf("%v", data) == "<nil>" {
		data = gin.H{}
	}
	if detail == nil {
		detail = gin.H{}
	}
	resp := Response{
		Code:   strconv.Itoa(msgCode),
		Msg:    msg,
		Detail: detail,
		Data:   data,
	}
	ctx.JSON(httpCode, resp)
}

func getMsg(ctx context.Context, msgCode int) string {
	msg := ""
	switch msgCode {
	case code.ErrDataNotExist.Code():
		msg = localize.Translate(ctx, "以下数据不存在：")
	case code.ErrNetworkConflict.Code():
		msg = localize.Translate(ctx, "以下数据存在冲突：")
	case code.ErrNetworkIllegal.Code():
		msg = localize.Translate(ctx, "以下数据非法：")
	case code.ErrExistRunningTask.Code():
		msg = localize.Translate(ctx, "以下数据存在正在执行的任务：")
	}

	return msg
}

// ExtendAttrsEntity 扩展属性
type ExtendAttrsEntity struct {
	ExtendAttrs map[string]interface{} `json:"extendAttrs" binding:"omitempty"`
}
