package redis

// Config Redis配置
type Config struct {
	Network           string `json:"network,omitempty"`
	Addr              string `json:"addr,omitempty"`
	ConnectionTimeout int64  `json:"connection_timeout,omitempty"`
	ReadTimeout       int64  `json:"read_timeout,omitempty"`
	WriteTimeout      int64  `json:"write_timeout,omitempty"`
	Password          string `json:"password,omitempty"`
	DB                int32  `json:"DB,omitempty"`
	MaxIdle           int32  `json:"max_idle,omitempty"`
	MaxActive         int32  `json:"max_active,omitempty"`
	TestOnBorrow      bool   `json:"test_on_borrow,omitempty"`
	IdleTimeout       int32  `json:"idle_timeout,omitempty"`
	Wait              bool   `json:"wait,omitempty"`
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

func (x *Config) GetConnectionTimeout() int64 {
	if x != nil {
		return x.ConnectionTimeout
	}
	return 0
}

func (x *Config) GetReadTimeout() int64 {
	if x != nil {
		return x.ReadTimeout
	}
	return 0
}

func (x *Config) GetWriteTimeout() int64 {
	if x != nil {
		return x.WriteTimeout
	}
	return 0
}

func (x *Config) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Config) GetDB() int32 {
	if x != nil {
		return x.DB
	}
	return 0
}

func (x *Config) GetMaxIdle() int32 {
	if x != nil {
		return x.MaxIdle
	}
	return 0
}

func (x *Config) GetMaxActive() int32 {
	if x != nil {
		return x.MaxActive
	}
	return 0
}

func (x *Config) GetTestOnBorrow() bool {
	if x != nil {
		return x.TestOnBorrow
	}
	return false
}

func (x *Config) GetIdleTimeout() int32 {
	if x != nil {
		return x.IdleTimeout
	}
	return 0
}

func (x *Config) GetWait() bool {
	if x != nil {
		return x.Wait
	}
	return false
}
