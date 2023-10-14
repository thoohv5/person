package database

// Config 数据库配置
type Config struct {
	// 数据库驱动
	Driver string `json:"driver,omitempty"`
	// 数据库连接
	Source                 string `json:"source,omitempty"`
	ConnMaxLifetimeSeconds int64  `json:"conn_max_lifetime_seconds,omitempty"`
	MaxOpenConns           int32  `json:"max_open_conns,omitempty"`
	MinIdleConns           int32  `json:"min_idle_conns,omitempty"`
	MaxRetries             int32  `json:"max_retries,omitempty"`
}

func (x *Config) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *Config) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Config) GetConnMaxLifetimeSeconds() int64 {
	if x != nil {
		return x.ConnMaxLifetimeSeconds
	}
	return 0
}

func (x *Config) GetMaxOpenConns() int32 {
	if x != nil {
		return x.MaxOpenConns
	}
	return 0
}

func (x *Config) GetMinIdleConns() int32 {
	if x != nil {
		return x.MinIdleConns
	}
	return 0
}

func (x *Config) GetMaxRetries() int32 {
	if x != nil {
		return x.MaxRetries
	}
	return 0
}
