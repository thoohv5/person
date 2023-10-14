package nats

import "encoding/json"

// IMessage 消息标准
type IMessage interface {
	// Marshal 序列化
	Marshal() ([]byte, error)
	// Unmarshal 反序列化
	Unmarshal(data []byte) error
}

type String struct {
	Data string
}

func (m *String) Marshal() ([]byte, error) {
	return []byte(m.Data), nil
}

func (m *String) Unmarshal(data []byte) error {
	m.Data = string(data)
	return nil
}

type Bytes struct {
	Data []byte
}

func (m *Bytes) Marshal() ([]byte, error) {
	return m.Data, nil
}

func (m *Bytes) Unmarshal(data []byte) error {
	m.Data = data
	return nil
}

type Map struct {
	Data map[string]interface{}
}

func (m *Map) Marshal() ([]byte, error) {
	return json.Marshal(m.Data)
}

func (m *Map) Unmarshal(data []byte) error {
	if m.Data == nil {
		m.Data = make(map[string]interface{})
	}
	return json.Unmarshal(data, m.Data)
}
