package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pg/v10/types"

	"github.com/thoohv5/person/internal/constant"
	ic "github.com/thoohv5/person/internal/context"
	"github.com/thoohv5/person/internal/provide/http"
	"github.com/thoohv5/person/internal/util"
)

type IQuery interface {
	Build(ctx context.Context, opts ...Option) (*Query, *func(err error) error, error)
	Where(opts ...QueryOption) IQuery
	GetQueryOption() []QueryOption
}

type Query struct {
	updateZero bool
	*pg.Query
	Q func(flag bool) (*pg.Query, error)
	*Result
}

func (q *Query) Insert(values ...interface{}) (pg.Result, error) {
	query, err := q.Q(true)
	if err != nil {
		return nil, err
	}
	return query.Insert(values...)
}

func (q *Query) Update(scan ...interface{}) (pg.Result, error) {
	query, err := q.Q(true)
	if err != nil {
		return nil, err
	}
	if q.updateZero {
		return query.Update(scan...)
	}
	return query.UpdateNotZero(scan...)
}

func (q *Query) Delete(values ...interface{}) (pg.Result, error) {
	query, err := q.Q(true)
	if err != nil {
		return nil, err
	}
	return query.Delete(values...)
}

func Create(Q func(flag bool) (*pg.Query, error), query *pg.Query, result *Result, updateZero bool) *Query {
	return &Query{
		Q:          Q,
		Query:      query,
		Result:     result,
		updateZero: updateZero,
	}
}

// condition 条件
type condition struct {
	// 参数
	opts []QueryOption
}

func New() IQuery {
	return &condition{}
}

func Where(opts ...QueryOption) IQuery {
	return &condition{
		opts: opts,
	}
}

// Build 构建
func (c *condition) Build(ctx context.Context, opts ...Option) (*Query, *func(err error) error, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	// if len(o.Models) == 0 {
	// 	return nil, nil, errors.WithLM("invalid model")
	// }

	call := func(err error) error {
		return err
	}
	var query *pg.Query
	if o.TableExpr != nil {
		query = o.TableExpr(o.Db, o.Db.ModelContext(ctx).Clone())
	} else if len(o.Models) > 0 {
		query = o.Db.ModelContext(ctx, o.Models...).Clone()
	} else {
		return nil, nil, errors.New("invalid model")
	}

	if o.QueryFunc != nil {
		o.QueryFunc(query)
	}

	for _, opt := range c.opts {
		opt(query)
	}

	return Create(func(flag bool) (*pg.Query, error) {
		if _, ok := o.Db.(*pg.Tx); !ok && flag {
			tx, err := o.Db.(*pg.DB).BeginContext(ctx)
			if err != nil {
				return nil, err
			}
			call = func(err error) error {
				if err != nil {
					rErr := tx.RollbackContext(ctx)
					if rErr != nil {
						err = fmt.Errorf("err: %w, rErr: %v", err, rErr)
					}
					return err
				}
				return tx.CommitContext(ctx)
			}
			if err = Message(ctx, tx); err != nil {
				return nil, err
			}
			return query.DB(tx), nil
		}
		return query, nil
	}, query, o.Result, o.updateZero), &call, nil
}

// Where 携带条件
func (c *condition) Where(opts ...QueryOption) IQuery {
	c.opts = append(c.opts, opts...)
	return c
}

// GetQueryOption 查询条件
func (c *condition) GetQueryOption() []QueryOption {
	return c.opts
}

// GetCommonQuery 获取通用查询条件
func GetCommonQuery(params http.BaseRequest, opts ...SearchOption) QueryOption {
	so := &searchOption{}
	for _, opt := range opts {
		opt(so)
	}
	// 查询条件
	limit := params.PageSize
	start := (params.PageNum - 1) * params.PageSize
	sorts := make([]string, 0)
	if sf := params.SortField; sf != "" {
		field := so.deal(util.Camel2Underline(sf))
		fs := strings.Split(field, ".")
		if len(fs) == 2 {
			if _, ok := Check(fs[0], fs[1]); ok || util.IsCustomField(fs[1]) {
				sorts = []string{fmt.Sprintf("%s %s", field, strings.ToLower(params.SortOrder))}
			}
		}
	}

	return Common(BaseRequest{
		Limit: limit,
		Start: start,
		Sorts: sorts,
	})
}

type QueryOption func(query *pg.Query)

func Common(br BaseRequest) QueryOption {
	return func(query *pg.Query) {
		if fields := br.Fields; len(fields) > 0 {
			query = query.Column(fields...)
		}

		if start := br.Start; start > 0 {
			query = query.Offset(int(start))
		}

		if limit := br.Limit; limit > 0 {
			query = query.Limit(int(limit))
		}

		if sorts := br.Sorts; len(sorts) > 0 {
			for _, sort := range sorts {
				ind := strings.Index(sort, " ")
				if ind != -1 {
					field := util.Camel2Underline(sort[:ind])
					sort1 := sort[ind+1:]
					fieldArr := strings.Split(field, ".")
					if len(fieldArr) == 2 {
						if fa := fieldArr[1]; util.IsCustomField(fa) {
							fieldArr[1] = fmt.Sprintf("extend_attrs->'extend_attrs'->>'%s'", fa)
							field = strings.Join(fieldArr, ".")
						}
					}
					query = query.OrderExpr(fmt.Sprintf("%s ?", field), types.Safe(sort1))
				} else if util.IsCustomField(sort) {
					sort = fmt.Sprintf("extend_attrs->'extend_attrs'->>'%s'", sort)
					query = query.OrderExpr(util.Camel2Underline(sort))
				} else {
					query = query.Order(util.Camel2Underline(sort))
				}
			}
		}
	}
}

func PgLogicalEmitMessage(ctx context.Context, tx *pg.Tx, key, val string) error {
	if _, err := tx.Exec("SELECT pg_logical_emit_message(true, ?, ?);", key, val); err != nil {
		if err = tx.RollbackContext(ctx); err != nil {
			return err
		}
		return err
	}
	return nil
}

func Message(ctx context.Context, tx *pg.Tx) error {
	keys := make([]string, 0)
	values := make(map[string]interface{})
	for _, lm := range ic.GetLogicalMessage() {
		if icx, ok := ctx.(*ic.Context); ok && lm.Key() == icx.Key() {
			lm = icx.ILogicalMessage
		} else if err := ic.GetMessage(ctx, lm); err != nil {
			if err.Error() == ic.ErrNotFoundKey {
				err = nil
				continue
			}
			return err
		}
		keys = append(keys, lm.Key())
		values[lm.Key()] = lm
	}
	// 没有基础base加一个
	if _, ok := values[ic.Key()]; !ok {
		blm := ic.FromContextBase(ctx)
		keys = append(keys, blm.Key())
		values[blm.Key()] = blm
	}
	marshal, err := json.Marshal(values)
	if err != nil {
		return err
	}
	return PgLogicalEmitMessage(ctx, tx, strings.Join(keys, ","), string(marshal))
}

type searchOption struct {
	deal func(tag string) string
}

// SearchOption 搜索可选字段
type SearchOption func(s *searchOption)

// WithDeal 可选字段处理
func WithDeal(deal func(tag string) string) SearchOption {
	return func(s *searchOption) {
		s.deal = deal
	}
}

// DefaultDeal 默认Deal
func DefaultDeal(prefix string) SearchOption {
	defaultDeal := func(tag string) string {
		return fmt.Sprintf("%s.%s", prefix, tag)
	}
	return func(s *searchOption) {
		s.deal = defaultDeal
	}
}

// DealSpecialSearch 处理特殊搜索
func DealSpecialSearch(query *pg.Query, search http.Search, opts ...SearchOption) {
	so := &searchOption{}
	for _, opt := range opts {
		opt(so)
	}
	switch search.Option {
	case constant.SearchOptionEqual:
		query.Where(fmt.Sprintf("%s = ?", search.Key), search.Value)
	case constant.SearchOptionLTE:
		query.Where(fmt.Sprintf("%s <= ?", search.Key), search.Value)
	case constant.SearchOptionGTE:
		query.Where(fmt.Sprintf("%s >= ?", search.Key), search.Value)
	case constant.SearchOptionLT:
		query.Where(fmt.Sprintf("%s < ?", search.Key), search.Value)
	case constant.SearchOptionGT:
		query.Where(fmt.Sprintf("%s > ?", search.Key), search.Value)
	}
}

// DealSearch 处理搜索
func DealSearch(query *pg.Query, search http.Search, opts ...SearchOption) {
	/**
	1. 搜索的字段是否带有表前缀
	2. 搜索的字段是否是自定义属性搜索
	3. 字段是否支持聚合OR
	*/

	so := &searchOption{}
	for _, opt := range opts {
		opt(so)
	}

	// 小驼峰转下划线
	search.Key = util.Camel2Underline(search.Key)

	tableName := ""
	// 可能字段有表名称
	if tIdx := strings.Index(search.Key, "."); tIdx != -1 {
		tableName = search.Key[:tIdx]
		search.Key = search.Key[tIdx+1:]
	}
	// 可能多个字段，聚合搜索
	fields := strings.Split(search.Key, ",")
	for idx, item := range fields {
		if so.deal != nil {
			item = so.deal(item)
			// 可能字段有表名称
			if tIdx := strings.Index(item, "."); tIdx != -1 {
				tableName = item[:tIdx]
				item = item[tIdx+1:]
			}
		}
		// 自定义字段处理一下
		if util.IsCustomField(item) {
			// 检查 FIX
			if search.Option == constant.SearchOptionIn {
				query.Where(fmt.Sprintf("%s.extend_attrs @> '{\"extend_attrs\": {\"%s\": [\"%s\"]}}'", tableName, item, strings.Join(strings.Split(fmt.Sprintf("%v", search.Value), "|"), "\",\"")))
				return
			}
			item = fmt.Sprintf("extend_attrs->'extend_attrs'->>'%s'", item)
			search.Value = fmt.Sprintf("%v", search.Value)
		}
		if len(tableName) > 0 {
			a, ok := Check(tableName, item)
			if !ok {
				fields[idx] = fmt.Sprintf("%s.%s", tableName, item)
				continue
			}
			switch a.Type {
			case reflect.TypeOf((*net.HardwareAddr)(nil)).Elem():
				search.Value = strings.Join(strings.FieldsFunc(fmt.Sprintf("%v", search.Value), func(r rune) bool {
					return r == ':' || r == '-'
				}), "")
				fields[idx] = fmt.Sprintf("REPLACE(%s.%s::text, ':', '')", tableName, item)
			case reflect.TypeOf((*net.Addr)(nil)).Elem():
				fields[idx] = fmt.Sprintf("host(%s.%s)", tableName, item)
			default:
				fields[idx] = fmt.Sprintf("%s.%s", tableName, item)
			}
		}
	}

	// 防撞处理
	rv := reflect.ValueOf(search.Value)
	switch rv.Kind() {
	case reflect.String:
		if search.Value.(string) == "" {
			// 空字符串不处理
			return
		}
	}

	// 数组组装
	query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
		format, param := GetFormat(search.Option, search.Value)
		for _, field := range fields {
			query = query.WhereOr(fmt.Sprintf(format, field), param)
		}
		return query, nil
	})
}

func GetFormat(option string, value interface{}) (format string, param interface{}) {
	format = ""
	param = value
	switch option {
	// 模糊匹配
	case constant.SearchOptionLike:
		format = "%s ilike ?"
		if vs, ok := value.(string); ok {
			param = fmt.Sprintf("%%%v%%", escapeLikePattern(vs))
		} else {
			param = fmt.Sprintf("%%%v%%", value)
		}
	// 开头匹配
	case constant.SearchOptionStartWith:
		format = "%s ilike ?"
		if vs, ok := value.(string); ok {
			param = fmt.Sprintf("%v%%", escapeLikePattern(vs))
		} else {
			param = fmt.Sprintf("%v%%", value)
		}
	// 等于
	case constant.SearchOptionEqual:
		format = "%s = ?"
	// 大于
	case constant.SearchOptionGT:
		format = "%s > ?"
	// 小于
	case constant.SearchOptionLT:
		format = "%s < ?"
	// 大于等于
	case constant.SearchOptionGTE:
		format = "%s >= ?"
	// 小于等于
	case constant.SearchOptionLTE:
		format = "%s <= ?"
	// 在...之中
	case constant.SearchOptionIn:
		format = "%s in (?)"
		if vs, ok := value.(string); ok {
			param = pg.In(strings.Split(vs, ","))
		} else {
			param = pg.In(value)
		}
	// 包含
	case constant.SearchOptionInclude:
		format = "%s @> ?"
	default:
		format = "%s = ?"
	}
	return format, param
}

func escapeLikePattern(pattern string) string {
	pattern = escapeCharacter(pattern, "\\")
	pattern = escapeCharacter(pattern, "%")
	pattern = escapeCharacter(pattern, "_")
	return pattern
}

func escapeCharacter(text, char string) string {
	return strings.ReplaceAll(text, char, "\\"+char)
}
