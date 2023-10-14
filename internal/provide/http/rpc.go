// Package http RPC
package http

import (
	"context"
	nethttp "net/http"
	"strconv"
	"time"

	"github.com/thoohv5/person/internal/code"
	ic "github.com/thoohv5/person/internal/context"
	uh "github.com/thoohv5/person/internal/util/http"
	"github.com/thoohv5/person/internal/util/rpc"
)

const (
	// URL 服务地址
	URL = "https://127.0.0.1"
)

// Base RPC基础
type Base struct {
	opts *options
}

// InitBase 初始化
func InitBase(o ...Option) *Base {
	opts := &options{
		url:     URL,
		timeout: 3 * time.Second,
	}
	for i := 0; i < len(o); i++ {
		o[i](opts)
	}
	return &Base{
		opts: opts,
	}
}

type options struct {
	url     string
	timeout time.Duration
}

// Option 可选项目
type Option func(o *options)

// WithURL 指定URL
func WithURL(url string) Option {
	return func(o *options) {
		o.url = url
	}
}

// WithTimeout 指定超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// Url URL
func (b *Base) Url() string {
	return b.opts.url
}

// Timeout Timeout
func (b *Base) Timeout() time.Duration {
	return b.opts.timeout
}

// getRPC 获取RPC
func (b *Base) getRPC() rpc.IRpc {
	opts := make([]uh.Option, 0)
	// tlsConfig, err := InitTLSConfig()
	// if err == nil {
	// 	opts = append(opts, uh.WithTLSClientConfig(tlsConfig))
	// }
	return rpc.New(rpc.WithOption(opts...))
}

// 检查Resp
func (b *Base) checkResp(ret *Response, resp *nethttp.Response) error {
	if !code.Success.EqualStr(ret.Code) && len(ret.Code) > 0 {
		codeInt, err := strconv.Atoi(ret.Code)
		if err != nil {
			return err
		}
		return code.NewWithMessage(codeInt, ret.Msg)
	}
	if resp.StatusCode != nethttp.StatusOK {
		return code.ErrRequestFail
	}
	return nil
}

// 默认Opts
func (b *Base) defaultOpts(ctx context.Context, resp *nethttp.Response) []uh.Option {
	return []uh.Option{uh.WithHeader(ic.ToMap(ctx)), uh.WithResponse(resp)}
}

// Get GET
func (b *Base) Get(ctx context.Context, url string, param interface{}, data interface{}, op ...uh.Option) (err error) {
	resp := &nethttp.Response{}
	ret := &Response{
		Data: data,
	}
	// Get
	if err = b.getRPC().Get(ctx, url, param, ret, append(b.defaultOpts(ctx, resp), op...)...); err != nil {
		return err
	}
	return b.checkResp(ret, resp)
}

// Post POST
func (b *Base) Post(ctx context.Context, url string, param interface{}, data interface{}, op ...uh.Option) (err error) {
	resp := &nethttp.Response{}
	ret := &Response{
		Data: data,
	}
	// Post
	if err = b.getRPC().Post(ctx, url, param, ret, append(b.defaultOpts(ctx, resp), op...)...); err != nil {
		return err
	}
	return b.checkResp(ret, resp)
}

// Put PUT
func (b *Base) Put(ctx context.Context, url string, param interface{}, data interface{}, op ...uh.Option) (err error) {
	resp := &nethttp.Response{}
	ret := &Response{
		Data: data,
	}
	// Put
	if err = b.getRPC().Put(ctx, url, param, ret, append(b.defaultOpts(ctx, resp), op...)...); err != nil {
		return err
	}
	return b.checkResp(ret, resp)
}

// Delete DELETE
func (b *Base) Delete(ctx context.Context, url string, param interface{}, data interface{}, op ...uh.Option) (err error) {
	resp := &nethttp.Response{}
	ret := &Response{
		Data: data,
	}
	// Delete
	if err = b.getRPC().Delete(ctx, url, param, ret, append(b.defaultOpts(ctx, resp), op...)...); err != nil {
		return err
	}
	return b.checkResp(ret, resp)
}
