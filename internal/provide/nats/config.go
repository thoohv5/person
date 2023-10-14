package nats

// Config 基础配置
type Config struct {
	Url        string `json:"url,omitempty"`
	ClientName string `json:"client_name,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Token      string `json:"token,omitempty"`
	// 最大重连次数
	MaxReconnect int32 `json:"max_reconnect,omitempty"`
	// 重连间隔--单位秒
	ReconnectTimeWait int32 `json:"reconnect_time_wait,omitempty"`
	// 连接超时--单位秒
	ConnectTimeout int32 `json:"connect_timeout,omitempty"`
	// Stream
	Streams []*JetStream `json:"streams,omitempty"`
}

func (x *Config) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Config) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

func (x *Config) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Config) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Config) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Config) GetMaxReconnect() int32 {
	if x != nil {
		return x.MaxReconnect
	}
	return 0
}

func (x *Config) GetReconnectTimeWait() int32 {
	if x != nil {
		return x.ReconnectTimeWait
	}
	return 0
}

func (x *Config) GetConnectTimeout() int32 {
	if x != nil {
		return x.ConnectTimeout
	}
	return 0
}

func (x *Config) GetStreams() []*JetStream {
	if x != nil {
		return x.Streams
	}
	return nil
}

// JetStream 配置
type JetStream struct {
	// 名称
	Name string `json:"name,omitempty"`
	// 储存方式：0-file, 1-memory, default: 0
	Storage int32 `json:"storage,omitempty"`
}

func (x *JetStream) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *JetStream) GetStorage() int32 {
	if x != nil {
		return x.Storage
	}
	return 0
}

// Producer 配置
type Producer struct {
	// 订阅主题，格式：JetStream配置name.XXX
	Subject string `json:"subject,omitempty"`
}

func (x *Producer) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

// Consumer 配置
type Consumer struct {
	// 名称
	Name string `json:"name,omitempty"`
	// 订阅主题，格式：JetStream配置name.XXX
	Subject string `json:"subject,omitempty"`
	// 消息投递策略： 1-all, 2-last, 3-new, 4-by_start_time, 5-by_start_sequence, 6-last_per_subject, default: 3
	DeliverPolicy int32 `json:"deliver_policy,omitempty"`
	// 消息确认策略：1-none, 2-all, 3-explicit，default: 3
	AckPolicy int32 `json:"ack_policy,omitempty"`
	// 每次拉取的消息数量，default: 200
	Fetch int32 `json:"fetch,omitempty"`
}

func (x *Consumer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Consumer) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Consumer) GetDeliverPolicy() int32 {
	if x != nil {
		return x.DeliverPolicy
	}
	return 0
}

func (x *Consumer) GetAckPolicy() int32 {
	if x != nil {
		return x.AckPolicy
	}
	return 0
}

func (x *Consumer) GetFetch() int32 {
	if x != nil {
		return x.Fetch
	}
	return 0
}
