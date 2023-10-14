// Package util
package util

import (
	"time"

	"github.com/robfig/cron/v3"
)

// GetCronNextExecAt 获取cron下次执行时间
func GetCronNextExecAt(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schl, err := parser.Parse(s)
	if err != nil {
		return time.Now(), err
	}

	return time.Time(schl.Next(time.Now())), nil
}
