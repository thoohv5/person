package enum

import (
	"context"
	"time"

	"github.com/thoohv5/person/internal/localize"
)

// TimeUnit 时间单位
//
//go:generate stringer -type TimeUnit -linecomment
type TimeUnit int32

// Second .
const (
	Second TimeUnit = iota + 1 // 秒
	Minute                     // 分钟
	Hour                       // 小时
	Day                        // 天
)

// Text 文本
func (i TimeUnit) Text(ctx context.Context) string {
	return localize.Translate(ctx, i.String())
}

// Duration 持续时间
func (i TimeUnit) Duration() time.Duration {
	td := time.Second
	switch i {
	case Second:
		td = time.Second
	case Minute:
		td = time.Minute
	case Hour:
		td = time.Hour
	case Day:
		td = 24 * time.Hour
	}
	return td
}
