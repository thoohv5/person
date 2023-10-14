// Package config 配置
package config

import (
	"fmt"
	"io"
	"log"

	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/internal/util"

	pConfig "github.com/thoohv5/person/pkg/config"
	"github.com/thoohv5/person/pkg/config/env"
	"github.com/thoohv5/person/pkg/config/file"
)

type conf struct {
	IBootstrap
	c pConfig.Config
}

// Config 配置标准
type Config interface {
	IBootstrap
	Close() error
	GetConfig() pConfig.Config
}

const (
	_envPrefix  = "THOOH_"
	_fileSuffix = ".yaml"
)

// New 创建
func New(path string) (Config, error) {
	c := pConfig.New(
		pConfig.WithSource(
			file.NewSource(path, _fileSuffix),
			env.NewSource(_envPrefix),
		),
		pConfig.WithLogger(newDefaultLogger()),
	)

	//  解析
	if err := c.Load(); nil != err {
		return nil, fmt.Errorf("config load err, err:%w", err)
	}

	// 转换
	bs := new(Bootstrap)
	if err := c.Scan(bs); nil != err {
		return nil, fmt.Errorf("config scan err, err:%w", err)
	}

	bs.Path = path

	return &conf{
		IBootstrap: bs,
		c:          c,
	}, nil
}

func (c *conf) Close() error {
	return c.c.Close()
}

func (c *conf) GetConfig() pConfig.Config {
	return c.c
}

// GetLogger 日志配置
func (c *conf) GetLogger() *logger.Config {
	l := c.IBootstrap.GetLogger()
	if l == nil {
		return &logger.Config{
			Out:   "std",
			Level: "debug",
		}
	}
	if l.File == nil {
		l.File.Path = "../../logs"
	}
	l.File.Path = util.AbPath(l.GetFile().GetPath())
	return l
}

// NewMockConfig mock配置
func NewMockConfig(bootstrap IBootstrap) Config {
	return &conf{
		IBootstrap: bootstrap,
	}
}

type defaultLogger struct {
	*log.Logger
}

func newDefaultLogger() pConfig.ILogger {
	return &defaultLogger{
		Logger: log.New(io.Discard, "config ", log.Lshortfile|log.Lmicroseconds|log.Ldate),
	}
}

func (d *defaultLogger) Debugf(msg string, values ...interface{}) {
	d.Printf(msg, values...)
}

func (d *defaultLogger) Infof(msg string, values ...interface{}) {
	d.Printf(msg, values...)
}

func (d *defaultLogger) Warnf(msg string, values ...interface{}) {
	d.Printf(msg, values...)
}

func (d *defaultLogger) Errorf(msg string, values ...interface{}) {
	d.Printf(msg, values...)
}
