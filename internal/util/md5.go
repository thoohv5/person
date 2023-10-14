// Package util
package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

// Md5 32位小写md5
func Md5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// SecretParam 密钥参数
type SecretParam struct {
	// Type 类型
	Type int32
	// Time 时间
	Time time.Time
	// Cron cron表达式
	Cron string
	// URL 地址
	URL string
	// Params 参数
	Params map[string]interface{}
}

// GenSecret 生成secret
// 任务中心密钥规则为：md5(type+exec_at+cron+exec_url+exec_params) 小写32位
func GenSecret(param SecretParam) (string, error) {
	bt := []byte(fmt.Sprintf("%d%s%s%s", param.Type, param.Time.String(), param.Cron, param.URL))
	t, err := json.Marshal(param.Params)
	if err != nil {
		return "", err
	}

	// 将params都序列化成map
	paramsMap := map[string]interface{}{}
	if uErr := json.Unmarshal(t, &paramsMap); uErr != nil {
		return "", uErr
	}
	t, err = json.Marshal(paramsMap)
	if err != nil {
		return "", err
	}
	bt = append(bt, t...)
	return Md5(string(bt)), nil
}
