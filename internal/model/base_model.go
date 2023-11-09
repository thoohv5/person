package model

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/samber/lo"

	"github.com/thoohv5/person/internal/util"
)

type BaseModel struct {
	IsIgnoreUpdate []string  `pg:"-"`
	ID             int32     `json:"id" pg:"id,pk"`
	CreatedTime    time.Time `json:"createdTime" pg:"created_time,notnull,default:now(),type:timestamp with time zone"`
	UpdatedTime    time.Time `json:"updatedTime" pg:"updated_time,notnull,default:now(),type:timestamp with time zone"`
}

type ExtendAttrsModel struct {
	BaseModel
	// 拓展字段：目前设计为存储自定义字段，后续还需要加不重要或定制字段时，也可以使用它。目前格式：暂定
	// {
	//	"extend_attrs":
	//		{
	//			"key1": "value",
	//			"key2": "value2"
	//		}
	// }
	ExtendAttrs ExtendAttr `json:"extend_attrs,omitempty" pg:"extend_attrs,notnull,default:'{}',use_zero"`
}

func (m *ExtendAttrsModel) GetExtendAttr() map[string]interface{} {
	if ea := m.ExtendAttrs.ExtendAttrs; len(ea) == 0 {
		return make(map[string]interface{})
	}
	return m.ExtendAttrs.ExtendAttrs
}

// ExtendAttr 拓展字段具体实体
type ExtendAttr struct {
	// 自定义字段所用
	ExtendAttrs map[string]interface{} `json:"extend_attrs,omitempty" pg:"extend_attrs,notnull,default:'{}'"`
}

func (m ExtendAttr) IsZero() bool {
	return m.ExtendAttrs == nil || len(m.ExtendAttrs) == 0
}

var _ pg.BeforeInsertHook = (*BaseModel)(nil)

func (m *BaseModel) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if m.CreatedTime.IsZero() {
		m.CreatedTime = now
	}
	if m.UpdatedTime.IsZero() {
		m.UpdatedTime = now
	}
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*BaseModel)(nil)

func (m *BaseModel) BeforeUpdate(ctx context.Context) (context.Context, error) {
	if len(m.IsIgnoreUpdate) == 0 || !lo.Contains(m.IsIgnoreUpdate, "updated_time") {
		m.UpdatedTime = time.Now()
	}
	return ctx, nil
}

var _ pg.AfterSelectHook = (*ExtendAttrsModel)(nil)

func (m *ExtendAttrsModel) AfterSelect(ctx context.Context) error {
	if m.ExtendAttrs.ExtendAttrs == nil {
		m.ExtendAttrs.ExtendAttrs = make(map[string]interface{})
	}
	return nil
}

func (m *BaseModel) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *BaseModel) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

func IsNotErrNoRows(err error) bool {
	return err != nil && err != pg.ErrNoRows
}

func IsErrNoRows(err error) bool {
	return err != nil && err == pg.ErrNoRows
}

func IntegrityViolation(err error) error {
	if err == nil {
		return nil
	}
	if pe, ok := err.(pg.Error); ok && pe.IntegrityViolation() {
		re := regexp.MustCompile(`Key .*\((.+)\) conflicts with existing key.*\((.+)\)`)
		keys := re.FindStringSubmatch(pe.Field('D'))
		if len(keys) > 1 {
			return fmt.Errorf("数据冲突 %s", strings.Join(keys[1:], ","))
		}
		re2 := regexp.MustCompile(`Key .*\((.+)\) already exists`)
		keys2 := re2.FindStringSubmatch(pe.Field('D'))
		if len(keys2) > 1 {
			return fmt.Errorf("%s 数据已存在", strings.Join(keys2[1:], ","))
		}
	}
	return err
}

// GetTableName 获取表名
// GetTableName((*Alarm)(nil)) => alarm
func GetTableName(entity interface{}) string {
	rt := reflect.TypeOf(entity)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return strings.Trim((string)(orm.GetTable(rt).SQLName), `"`)
}

// InitExtendAttrsModel 初始化
func InitExtendAttrsModel(extendAttrs map[string]interface{}) ExtendAttrsModel {
	return ExtendAttrsModel{
		ExtendAttrs: ExtendAttr{
			ExtendAttrs: extendAttrs,
		},
	}
}

// Attribute 属性
type Attribute struct {
	// pg tag
	Tag string
	// 类型
	Type reflect.Type
}

// key: tableName => value: {key : field => attribute}
var mTableAttribute = make(map[string]map[string]*Attribute)

// Check 检查
func Check(tableName, tag string) (*Attribute, bool) {
	m, ok := mTableAttribute[tableName]
	if !ok {
		return nil, false
	}
	attribute, ok := m[tag]
	if !ok {
		return nil, false
	}
	return attribute, true
}

// Get 表前缀
func Get(tnList []string) func(tag string) string {
	if len(tnList) == 0 {
		return func(tag string) string {
			return tag
		}
	}
	data := make(map[string]string)
	for _, t := range tnList {
		for field := range mTableAttribute[t] {
			data[field] = t
		}
	}
	return func(tag string) string {
		tn := tnList[0]
		if t, ok := data[tag]; ok {
			tn = t
		}
		return fmt.Sprintf("%s.%s", tn, tag)
	}
}

// Register 注册
func Register(m interface{}) {
	tableName := GetTableName(m)
	if _, ok := mTableAttribute[tableName]; !ok {
		mTableAttribute[tableName] = make(map[string]*Attribute)
	}

	rt := reflect.TypeOf(m).Elem()
	splitStruct(rt, mTableAttribute[tableName])
}

func splitStruct(rt reflect.Type, m map[string]*Attribute) {
	for i := 0; i < rt.NumField(); i++ {
		rtf := rt.Field(i)
		if rtf.Anonymous {
			splitStruct(rtf.Type, m)
			continue
		}
		tag, ok := rtf.Tag.Lookup("pg")
		if tag == "-" {
			continue
		}
		if !ok {
			tag = util.Camel2Underline(rtf.Name)
		}
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		m[tag] = &Attribute{
			Tag:  tag,
			Type: rtf.Type,
		}
	}
}
