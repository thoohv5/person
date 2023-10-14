package code

import (
	"net/http"
)

var (
	// 全局Code
	allCode = map[int]ICode{}

	// Success 请求成功
	Success = &Code{status: http.StatusOK, code: 0, msg: "请求成功"}
)
