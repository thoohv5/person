// Package code 错误码
package code

import (
	"fmt"
)

// ICode 错误码标识
type ICode interface {
	// Status ,Status Code
	Status() int
	// Code ,Error Code
	Code() int
	// Error ,Error msg
	Error() string
}

// Code 错误码
type Code struct {
	status int
	code   int
	msg    string
	data   []interface{}
	detail interface{}
}

// Status 状态
func (c *Code) Status() int {
	return c.status
}

// Code 码
func (c *Code) Code() int {
	return c.code
}

// Error 错误
func (c *Code) Error() string {
	return c.msg
}

// Detail 详情
func (c *Code) Detail() interface{} {
	return c.detail
}

// CodeStr 将code转成string
func (c *Code) CodeStr() string {
	return fmt.Sprintf("%d", c.code)
}

// EqualStr 和str错误码比较
func (c *Code) EqualStr(strCode string) bool {
	return c.CodeStr() == strCode
}

// Equal 相等
func (c *Code) Equal(err error) bool {
	switch err.(type) {
	case *Code:
		return err.(*Code).Code() == c.Code()
	default:
		return c.Error() == err.Error()
	}
}

// Equal 相等
func Equal(code int, err error) bool {
	c := &Code{code: code}
	switch err.(type) {
	case *Code:
		return err.(*Code).Code() == c.Code()
	default:
		return c.Error() == err.Error()
	}
}

// New 创建错误码
func New(code int, data ...interface{}) ICode {
	return NewWithMessage(code, ErrCommon.Error(), data)
}

// NewWithMessage 创建错误码通过错误信息
func NewWithMessage(code int, msg string, data ...interface{}) (iCode ICode) {

	iCode, ok := allCode[code]
	if !ok {
		iCode = &Code{
			status: ErrCommon.Status(),
			code:   code,
			msg:    fmt.Sprintf(msg, data...),
		}
		// 不考虑竞争
		allCode[code] = iCode
	} else {
		iCode = &Code{
			status: iCode.Status(),
			code:   code,
			msg:    fmt.Sprintf(msg, data...),
		}
	}
	return
}

// WithData 携带动态数据
func WithData(iCode ICode, data ...interface{}) ICode {
	return &Code{
		status: iCode.Status(),
		code:   iCode.Code(),
		msg:    fmt.Sprintf(iCode.Error(), data...),
	}
}

// WithMessage 替换错误信息
func WithMessage(iCode ICode, msg string, data ...interface{}) ICode {
	return &Code{
		status: iCode.Status(),
		code:   iCode.Code(),
		msg:    fmt.Sprintf(msg, data...),
	}
}

// WithDetail 替换详细信息
func WithDetail(iCode ICode, detail interface{}) ICode {
	return &Code{
		status: iCode.Status(),
		code:   iCode.Code(),
		detail: detail,
	}
}
