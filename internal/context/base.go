package context

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"

	"github.com/thoohv5/person/internal/util"
	"github.com/thoohv5/person/internal/util/mapstructure"
)

type base struct {
	TraceID  string `header:"X-Request-Id" json:"traceID"`
	Lang     string `header:"X-Lang" json:"lang"`
	UserID   string `header:"X-User-Id" json:"userID"`
	ClientIP string `header:"X-Real-Ip" json:"realIP"`
	GroupIDs string `header:"X-Group-Ids" json:"groupIDs"`
}

const (
	_label      = "THOOH"
	_defaultTag = "header"
)

func init() {
	registerLogicalMessage(&base{})
}

type IBase interface {
	ILogicalMessage
	Additional(spanId string) IBase

	SetTraceID(traceID string) IBase
	GetTraceID() string
	SetLang(lang string) IBase
	GetLang() string
	SetUserID(userID string) IBase
	GetUserID() string
	SetClientIP(userID string) IBase
	GetClientIP() string
	SetGroupIDs(groupIDs map[string][]int32) IBase
	GetGroupIDs() map[string][]int32
}

var (
	_default = &base{}
)

func (m *base) Key() string {
	return fmt.Sprintf("%s-%s", _label, util.InitialLower(reflect.TypeOf(m).Elem().Name()))
}

func Key() string {
	return _default.Key()
}

func (m *base) Marshal(ctx context.Context) (string, error) {
	bs, err := json.Marshal(m)
	return string(bs), err
}

func (m *base) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), m)
}

func (m *base) Merge(message ILogicalMessage) {
	return
}

func (m *base) Additional(spanId string) IBase {
	m.TraceID = fmt.Sprintf("%s:%s", m.TraceID, spanId)
	return m
}

func WithBase(b IBase) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, b.Key(), b)
	}
}

func FromContextBase(ctx context.Context) IBase {
	value, ok := ctx.Value(Key()).(IBase)
	if !ok {
		value = &base{}
	}
	if traceID := value.GetTraceID(); len(traceID) == 0 {
		value.SetTraceID(GetDefaultTraceID())
	}
	if lang := value.GetLang(); len(lang) == 0 {
		value.SetLang(GetDefaultLang())
	}
	return value
}

func Copy(ctx context.Context) context.Context {
	b := FromContextBase(ctx)
	return WithBase(b)(context.Background())
}

func ToMap(ctx context.Context) map[string]string {
	ret := make(map[string]string)
	m := FromContextBase(ctx)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &ret,
		TagName:  _defaultTag,
	})
	if err != nil {
		panic(err)
	}
	if err = decoder.Decode(m); err != nil {
		panic(err)
	}
	return ret
}

func WithHeader() func(ctx context.Context, data map[string][]string) context.Context {
	return func(ctx context.Context, data map[string][]string) context.Context {
		b := &base{}
		rt := reflect.TypeOf(b).Elem()
		rv := reflect.Indirect(reflect.ValueOf(b))
		for i := 0; i < rt.NumField(); i++ {
			v, ok := data[GetTag(rt.Field(i))]
			if !ok || len(v) == 0 {
				continue
			}
			if rvf := rv.Field(i); rvf.CanSet() {
				rvf.SetString(v[0])
			}
		}
		return WithBase(b)(ctx)
	}
}

/*********clientIP***********/

func GetClientIPLabel() string {
	fieldByName, b := reflect.TypeOf(&base{}).Elem().FieldByName("ClientIP")
	if !b {
		return ""
	}
	return GetTag(fieldByName)
}

func (m *base) SetClientIP(clientIP string) IBase {
	m.ClientIP = clientIP
	return m
}

func (m *base) GetClientIP() string {
	return m.ClientIP
}

func WithClientIP(clientIP string) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		b := FromContextBase(ctx)
		return WithBase(b.SetClientIP(clientIP))(ctx)
	}
}

func FromCtxClientIP(ctx context.Context) string {
	return FromContextBase(ctx).GetClientIP()
}

/*********traceID***********/

func (m *base) SetTraceID(traceID string) IBase {
	m.TraceID = traceID
	return m
}

func (m *base) GetTraceID() string {
	return m.TraceID
}

func GetTraceIDLabel() string {
	fieldByName, b := reflect.TypeOf(&base{}).Elem().FieldByName("TraceID")
	if !b {
		return ""
	}
	return GetTag(fieldByName)
}

func FromCtxTraceID(ctx context.Context) string {
	return FromContextBase(ctx).GetTraceID()
}

/*********lang***********/

func GetLangLabel() string {
	fieldByName, b := reflect.TypeOf(&base{}).Elem().FieldByName("Lang")
	if !b {
		return ""
	}
	return GetTag(fieldByName)
}

func (m *base) SetLang(lang string) IBase {
	m.Lang = lang
	return m
}

func (m *base) GetLang() string {
	return m.Lang
}

func FromCtxLang(ctx context.Context) string {
	return FromContextBase(ctx).GetLang()
}

/*********userID***********/

func GetUserIDLabel() string {
	fieldByName, b := reflect.TypeOf(&base{}).Elem().FieldByName("UserID")
	if !b {
		return ""
	}
	return GetTag(fieldByName)
}

func (m *base) SetUserID(userID string) IBase {
	m.UserID = userID
	return m
}

func (m *base) GetUserID() string {
	return m.UserID
}

func WithUserID(userID string) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		b := FromContextBase(ctx)
		return WithBase(b.SetUserID(userID))(ctx)
	}
}

func WithCtxUserID(callback func(ctx context.Context, userID string) bool) func(ctx context.Context, userID string) (context.Context, bool) {
	return func(ctx context.Context, userID string) (context.Context, bool) {
		return WithUserID(userID)(ctx), callback(ctx, userID)
	}
}

func FromCtxUserID(ctx context.Context) string {
	return FromContextBase(ctx).GetUserID()
}

/*********GroupIDs***********/

func WithGroupIDs(groupIDs map[string][]int32) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		b := FromContextBase(ctx)
		return WithBase(b.SetGroupIDs(groupIDs))(ctx)
	}
}

func WithCtxGroupIDs() func(ctx context.Context, groupIDs map[string][]int32) context.Context {
	return func(ctx context.Context, groupIDs map[string][]int32) context.Context {
		return WithGroupIDs(groupIDs)(ctx)
	}
}

func (m *base) SetGroupIDs(groupIDs map[string][]int32) IBase {
	var s []byte
	s, err := json.Marshal(groupIDs)
	if err != nil {
		s = []byte("{}")
	}
	m.GroupIDs = string(s)
	return m
}

func (m *base) GetGroupIDs() map[string][]int32 {
	var groupIDs map[string][]int32
	if err := json.Unmarshal([]byte(m.GroupIDs), &groupIDs); err != nil {
		groupIDs = map[string][]int32{}
	}
	return groupIDs
}

func FromCtxGroupIDs(ctx context.Context) map[string][]int32 {
	return FromContextBase(ctx).GetGroupIDs()
}

/*********默认值***********/

const (
	defaultLang   = "zh-CN"
	defaultUserID = "admin"
)

func GetDefaultTraceID() string {
	return uuid.NewString()
}

// GetDefaultLang 默认语言
func GetDefaultLang() string {
	return defaultLang
}

// GetDefaultUserID 默认用户ID
func GetDefaultUserID() string {
	return defaultUserID
}

/*********工具***********/

// GetTag 获取json tag
func GetTag(sf reflect.StructField) string {
	tag, ok := sf.Tag.Lookup(_defaultTag)
	if !ok {
		return ""
	}
	cTag := tag
	if idx := strings.Index(tag, ","); idx != -1 {
		cTag = tag[:idx]
	}
	return cTag
}
