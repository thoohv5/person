package http

// Config 配置
type Config struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Timeout int32  `json:"timeout,omitempty"`
	// 模式 release/debug
	Model string `json:"model,omitempty"`
	// 开启swag
	EnableSwag bool `json:"enable_swag,omitempty"`
	// 开启pprof
	EnablePprof bool `json:"enable_pprof,omitempty"`
}

func (x *Config) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Config) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Config) GetTimeout() int32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *Config) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Config) GetEnableSwag() bool {
	if x != nil {
		return x.EnableSwag
	}
	return false
}

func (x *Config) GetEnablePprof() bool {
	if x != nil {
		return x.EnablePprof
	}
	return false
}
