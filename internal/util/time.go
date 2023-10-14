package util

import (
	"fmt"
	"time"
)

// CalculateTime 计算天
func CalculateTime(start, end time.Time) string {
	// 计算两个时间之间的差异
	duration := end.Sub(start)
	// 将差异转换为天数
	days := int32(duration.Hours() / 24)
	hours := int32(duration.Hours()) % 24
	return fmt.Sprintf("%d天%d小时", days, hours)
}
