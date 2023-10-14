package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/thoohv5/person/pkg/log"

	"github.com/thoohv5/person/internal/util"
)

var _ log.Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	log *zap.Logger
	ctx context.Context
}

const (
	OutSep = ","
)

func New(loggerConf *Config) (log.Logger, func(), error) {

	var zlg *zap.Logger
	if loggerConf.Model == 1 {
		zlg = zap.New(getWriter(loggerConf))
	} else {
		zlg = zap.New(getWriter(loggerConf), zap.AddCaller(), zap.Development(), zap.AddCallerSkip(1))
	}

	zlg = zlg.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return &defaultLogger{log: zlg}, func() {
	}, nil
}

func pre(ctx context.Context, fields ...log.Field) []zap.Field {
	fs := NewFields()
	FieldContext(ctx)(fs)
	for _, field := range fields {
		field(fs)
	}

	items := make([]zap.Field, 0, len(fs.Data()))
	log.SMap(fs.Data(), func(key string, item *log.Entity) {
		switch item.Type {
		case reflect.String:
			val, ok := item.Value.(string)
			if !ok {
				val = fmt.Sprint(val)
			}
			items = append(items, zap.String(key, val))
		default:
			items = append(items, zap.Any(key, item.Value))
		}
	})
	return items
}

func (l *defaultLogger) AddCallerSkip(skip int) log.Logger {
	copy := *l
	copy.log = copy.log.WithOptions(zap.AddCallerSkip(skip))
	return &copy
}

func (l *defaultLogger) Debugc(ctx context.Context, msg string, fields ...log.Field) {
	l.log.Debug(msg, pre(ctx, fields...)...)
}
func (l *defaultLogger) Infoc(ctx context.Context, msg string, fields ...log.Field) {
	l.log.Info(msg, pre(ctx, fields...)...)
}
func (l *defaultLogger) Warnc(ctx context.Context, msg string, fields ...log.Field) {
	l.log.Warn(msg, pre(ctx, fields...)...)
}
func (l *defaultLogger) Errorc(ctx context.Context, msg string, fields ...log.Field) {
	l.log.Error(msg, pre(ctx, fields...)...)
}

func (l *defaultLogger) Debugf(ctx context.Context, msg string, values ...interface{}) {
	l.log.Debug(fmt.Sprintf(msg, values...), pre(ctx)...)
}
func (l *defaultLogger) Infof(ctx context.Context, msg string, values ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, values...), pre(ctx)...)
}
func (l *defaultLogger) Warnf(ctx context.Context, msg string, values ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, values...), pre(ctx)...)
}
func (l *defaultLogger) Errorf(ctx context.Context, msg string, values ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, values...), pre(ctx)...)
}

func (l *defaultLogger) Sync() error {
	// return l.log.Sync()
	return nil
}

func (l *defaultLogger) Close() error {
	return l.Sync()
}

type devNullLog struct {
	w io.Writer
}

func NewDevNullLog(w io.Writer) log.Logger {
	return &devNullLog{
		w: w,
	}
}

func (d *devNullLog) AddCallerSkip(skip int) log.Logger {
	return d
}

func (d *devNullLog) Debugc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Infoc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Warnc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}
func (d *devNullLog) Errorc(ctx context.Context, msg string, fields ...log.Field) {
	fmt.Fprintf(d.w, msg, pre(ctx, fields...))
}

func (d *devNullLog) Debugf(ctx context.Context, msg string, values ...interface{}) {
	fmt.Fprintf(d.w, msg, values...)
}
func (d *devNullLog) Infof(ctx context.Context, msg string, values ...interface{}) {
	fmt.Fprintf(d.w, msg, values...)
}
func (d *devNullLog) Warnf(ctx context.Context, msg string, values ...interface{}) {
	fmt.Fprintf(d.w, msg, values...)
}
func (d *devNullLog) Errorf(ctx context.Context, msg string, values ...interface{}) {
	fmt.Fprintf(d.w, msg, values...)
}

func getWriter(conf *Config) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "func",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	ws := make([]zapcore.WriteSyncer, 0)
	for _, out := range strings.Split(conf.GetOut(), OutSep) {
		switch out {
		case "std":
			ws = append(ws, zapcore.AddSync(os.Stdout))
		case "discard":
			ws = append(ws, zapcore.AddSync(io.Discard))
		default:
			ws = append(ws, zapcore.AddSync(getFileWriter(conf.GetFile())))
		}
	}
	var encoder zapcore.Encoder
	if conf.Type == "text" {
		encoder = NewTextEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(ws...),
		zap.NewAtomicLevelAt(parseLevel(conf.GetLevel())),
	)

}

// 日志类别: debug, warn, info，error
func parseLevel(level string) zapcore.Level {
	zl := zapcore.DebugLevel
	switch level {
	case "debug":
		zl = zapcore.DebugLevel
	case "warn":
		zl = zapcore.WarnLevel
	case "info":
		zl = zapcore.InfoLevel
	case "error":
		zl = zapcore.ErrorLevel
	}
	return zl
}

func getFileWriter(fc *File) io.Writer {
	if strings.HasPrefix(fc.GetPath(), ".") {
		fc.Path = util.AbPath(fc.GetPath())
	}
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", fc.GetPath(), fc.GetFileName()),
		MaxSize:    int(fc.GetMaxSize()),
		MaxBackups: int(fc.GetMaxBackups()),
		MaxAge:     int(fc.GetMaxAge()),
		Compress:   fc.GetCompress(),
	}
}
