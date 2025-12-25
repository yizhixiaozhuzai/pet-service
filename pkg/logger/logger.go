package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log  *zap.Logger
	slog *zap.SugaredLogger
	once sync.Once
)

// LogLevel 日志级别类型
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

// Init 初始化日志系统
func Init(level LogLevel, maxSize, maxBackups, maxAge int, compress bool, outputPath string) error {
	var initErr error
	once.Do(func() {
		initErr = initLogger(level, maxSize, maxBackups, maxAge, compress, outputPath)
	})
	return initErr
}

func initLogger(level LogLevel, maxSize, maxBackups, maxAge int, compress bool, outputPath string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 配置日志轮转
	writer := &lumberjack.Logger{
		Filename:   outputPath,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups, // 保留的旧日志文件最大数量
		MaxAge:     maxAge,     // 保留天数
		Compress:   compress,   // 压缩
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 解析日志级别
	var zapLevel zapcore.Level
	switch level {
	case LevelDebug:
		zapLevel = zapcore.DebugLevel
	case LevelInfo:
		zapLevel = zapcore.InfoLevel
	case LevelWarn:
		zapLevel = zapcore.WarnLevel
	case LevelError:
		zapLevel = zapcore.ErrorLevel
	case LevelFatal:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		zapLevel,
	)

	// 创建Logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	slog = log.Sugar()

	return nil
}

// Info 记录info级别日志
func Info(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[INFO] %s\n", msg)
		return
	}
	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Infow(msg, toZapFields(fs...)...)
}

// Debug 记录debug级别日志
func Debug(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[DEBUG] %s\n", msg)
		return
	}
	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Debugw(msg, toZapFields(fs...)...)
}

// Warn 记录warn级别日志
func Warn(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[WARN] %s\n", msg)
		return
	}
	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Warnw(msg, toZapFields(fs...)...)
}

// Error 记录error级别日志
func Error(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[ERROR] %s\n", msg)
		return
	}
	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Errorw(msg, toZapFields(fs...)...)
}

// Fatal 记录fatal级别日志
func Fatal(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[FATAL] %s\n", msg)
		os.Exit(1)
	}
	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Fatalw(msg, toZapFields(fs...)...)
}

// WithCaller 记录调用者信息的日志
func WithCaller(ctx context.Context, msg string, fields ...Field) {
	if slog == nil {
		fmt.Printf("[INFO] %s\n", msg)
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fields = append(fields, String("caller", fmt.Sprintf("%s:%d", filepath.Base(file), line)))
	}

	fs := append([]Field{traceIDFromContext(ctx)}, fields...)
	slog.Infow(msg, toZapFields(fs...)...)
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// String 创建string字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int 创建int字段
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 创建int64字段
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Any 创建any字段
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Error 创建error字段
func ErrorField(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func toZapFields(fields ...Field) []interface{} {
	fs := make([]interface{}, 0, len(fields)*2)
	for _, f := range fields {
		fs = append(fs, f.Key, f.Value)
	}
	return fs
}

// traceIDFromContext 从context中获取traceID
func traceIDFromContext(ctx context.Context) Field {
	if ctx == nil {
		return String("trace_id", "unknown")
	}
	// 这里可以根据实际使用的链路追踪库获取traceID
	// 例如: traceID, _ := tracer.TraceIDFromContext(ctx)
	return String("trace_id", "unknown")
}

// GetLogger 获取原生logger
func GetLogger() *zap.Logger {
	return log
}

// GetSugaredLogger 获取SugaredLogger
func GetSugaredLogger() *zap.SugaredLogger {
	return slog
}

// Sync 同步日志缓冲区
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
