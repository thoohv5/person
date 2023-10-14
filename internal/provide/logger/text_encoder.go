package logger

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type textEncoder struct {
	*zapcore.EncoderConfig
	*zapcore.MapObjectEncoder
}

var (
	_pool = buffer.NewPool()
)

func NewTextEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	if config.ConsoleSeparator == "" {
		// Use a default delimiter of '\t' for backwards compatibility
		config.ConsoleSeparator = "\t"
	}
	return textEncoder{
		EncoderConfig:    &config,
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

func (c textEncoder) Clone() zapcore.Encoder {
	return textEncoder{
		EncoderConfig:    c.EncoderConfig,
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

func (c textEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := _pool.Get()
	arr := &sliceArrayEncoder{}
	if c.TimeKey != "" && c.EncodeTime != nil {
		c.EncodeTime(ent.Time, arr)
	}

	// 从上下文中获取traceid，并移除这个字段
	fields = c.addTrace(arr, fields)

	if c.LevelKey != "" && c.EncodeLevel != nil {
		c.EncodeLevel(ent.Level, arr)
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		nameEncoder := c.EncodeName

		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, arr)
	}
	for i := range arr.elems {
		if i > 0 {
			line.AppendString(c.ConsoleSeparator)
		}
		fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	// Add the message itself.
	if c.MessageKey != "" {
		c.addSeparatorIfNecessary(line)
		line.AppendByte('"')
		line.AppendString(ent.Message)
		line.AppendByte('"')
	}

	c.writeContext(line, fields)

	if ent.Caller.Defined {
		if c.CallerKey != "" && c.EncodeCaller != nil {
			c.addSeparatorIfNecessary(line)
			fmt.Fprint(line, ent.Caller)
		}
		if c.FunctionKey != "" {
			c.addSeparatorIfNecessary(line)
			if ind := strings.LastIndexByte(ent.Caller.Function, '/'); ind != -1 {
				line.AppendString(ent.Caller.Function[ind+1:])
			} else {
				line.AppendString(ent.Caller.Function)
			}
		}
	}

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	line.AppendString(c.LineEnding)
	return line, nil
}

func (c textEncoder) addTrace(arr *sliceArrayEncoder, fields []zapcore.Field) []zapcore.Field {
	len1 := len(fields)
	if len1 > 0 {
		var traceId string
		for key, field := range fields {
			if field.Key == "trace_id" {
				traceId = field.String
				fields = append(fields[:key], fields[key+1:]...)
				break
			}
		}
		if traceId != "" {
			arr.AppendString(traceId)
		}
	}
	return fields
}

func (c textEncoder) writeContext(line *buffer.Buffer, fields []zapcore.Field) {
	data := zapcore.NewMapObjectEncoder()
	for i := range fields {
		fields[i].AddTo(data)
	}
	c.addSeparatorIfNecessary(line)

	statement, ok := data.Fields["statement"]
	if ok {
		fmt.Printf("%c[%d;%d;%dm SQL: %s%c[0m\n", 0x1B, 4, 0, 32, statement, 0x1B)
	}

	if len(data.Fields) > 0 {
		msg := make([]string, 0, len(data.Fields))
		for k, v := range data.Fields {
			if v == nil {
				continue
			}
			if rt := reflect.TypeOf(v); rt.Kind() == reflect.Ptr {
				v = reflect.ValueOf(v).Elem()
			}
			msg = append(msg, fmt.Sprintf("%s:%v", k, v))
		}
		_, err := line.WriteString(strings.Join(msg, " "))
		if err != nil {
			fmt.Printf("write context error: %v\n", err)
		}
	}
}

func (c textEncoder) addSeparatorIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendString(c.ConsoleSeparator)
	}
}

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

// sliceArrayEncoder is an ArrayEncoder backed by a simple []interface{}. Like
// the MapObjectEncoder, it's not designed for production use.
type sliceArrayEncoder struct {
	elems []interface{}
}

func (s *sliceArrayEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	enc := &sliceArrayEncoder{}
	err := v.MarshalLogArray(enc)
	s.elems = append(s.elems, enc.elems)
	return err
}

func (s *sliceArrayEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	m := zapcore.NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	s.elems = append(s.elems, m.Fields)
	return err
}

func (s *sliceArrayEncoder) AppendReflected(v interface{}) error {
	s.elems = append(s.elems, v)
	return nil
}

func (s *sliceArrayEncoder) AppendBool(v bool)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendByteString(v []byte)      { s.elems = append(s.elems, string(v)) }
func (s *sliceArrayEncoder) AppendComplex128(v complex128)  { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendComplex64(v complex64)    { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendDuration(v time.Duration) { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat64(v float64)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendFloat32(v float32)        { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt(v int)                { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt64(v int64)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt32(v int32)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt16(v int16)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendInt8(v int8)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendString(v string)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendTime(v time.Time)         { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint(v uint)              { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint64(v uint64)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint32(v uint32)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint16(v uint16)          { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUint8(v uint8)            { s.elems = append(s.elems, v) }
func (s *sliceArrayEncoder) AppendUintptr(v uintptr)        { s.elems = append(s.elems, v) }
