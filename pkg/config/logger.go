package config

import (
	"log"
	"os"
)

// ILogger is a logger interface.
type ILogger interface {
	Debugf(msg string, values ...interface{})
	Infof(msg string, values ...interface{})
	Warnf(msg string, values ...interface{})
	Errorf(msg string, values ...interface{})
}

type defaultLogger struct {
	*log.Logger
}

func NewDefaultLogger() ILogger {
	return &defaultLogger{
		Logger: log.New(os.Stdout, "config ", log.Lshortfile|log.Lmicroseconds|log.Ldate),
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
