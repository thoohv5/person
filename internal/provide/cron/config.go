package cron

// Config 配置
type Config struct {
	Enable bool `json:"enable,omitempty"`
	Debug  bool `json:"debug,omitempty"`
}

func (x *Config) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *Config) GetDebug() bool {
	if x != nil {
		return x.Debug
	}
	return false
}

type Timer struct {
	Spec string `json:"spec,omitempty"`
}

func (x *Timer) GetSpec() string {
	if x != nil {
		return x.Spec
	}
	return ""
}
