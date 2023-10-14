package context

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/thoohv5/person/internal/util"
)

type comptroller struct {
	Data map[string]string `json:"data"`
}

func init() {
	registerLogicalMessage(&comptroller{})
}

type IComptroller interface {
	ILogicalMessage
	SetData(key, value string) IComptroller
	Get(tag string) string
}

func NewComptroller() IComptroller {
	return &comptroller{
		Data: make(map[string]string),
	}
}

func (m *comptroller) Key() string {
	return util.InitialLower(reflect.TypeOf(m).Elem().Name())
}

func (m *comptroller) Marshal(ctx context.Context) (string, error) {
	bs, err := json.Marshal(m)
	return string(bs), err
}

func (m *comptroller) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), &m.Data)
}

func (m *comptroller) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Data)
}

func (m *comptroller) Merge(message ILogicalMessage) {
	c, ok := message.(*comptroller)
	if !ok {
		return
	}
	for key, val := range c.Data {
		if _, ok = m.Data[key]; !ok {
			m.Data[key] = val
		}
	}
	return
}

func (m *comptroller) SetData(key, value string) IComptroller {
	m.Data[key] = value
	return m
}

func (m *comptroller) Get(tag string) string {
	return m.Data[tag]
}
