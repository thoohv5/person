// Package base 基类
package cron

// Base 定时任务基础类
type Base struct {
}

// Enable 是否启用
func (b *Base) Enable() bool {
	return false
}
