package log

import (
	"OneDisk/lib/definition"
	"OneDisk/lib/format"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var printer *zap.Logger

// Initialize
// 初始化日志记录器
func Initialize() {
	// 配置日志分片
	logger := &lumberjack.Logger{
		Filename: definition.PathLog,
		MaxSize:  10,   // megabytes
		MaxAge:   30,   // days
		Compress: true, // enabled by default
	}

	// 日志打印配置
	configTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05 000"))
	}
	configTagEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    configTagEncoder,
		EncodeTime:     configTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder}

	// 配置 zap 日志记录器
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(logger),
		zap.InfoLevel,
	)

	printer = zap.New(core)
}

// Info
// 输出日志
// tag: 日志标签
// msg: 日志内容
// fields: 日志字段
func Info(tag string, msg string, fields ...zap.Field) {
	printer.Info("["+tag+"] "+msg, fields...)
}

// Debug
// 输出日志
// tag: 日志标签
// msg: 日志内容
// fields: 日志字段
func Debug(tag string, msg string, fields ...zap.Field) {
	printer.Debug("["+tag+"] "+msg, fields...)
}

// Warming
// 输出日志
// tag: 日志标签
// msg: 日志内容
// fields: 日志字段
func Warming(tag string, msg string, fields ...zap.Field) {
	printer.Warn("["+tag+"] "+msg, fields...)
}

// Error
// 输出错误日志
// tag: 日志标签
// msg: 日志内容
// fields: 日志字段
func Error(tag string, msg string, fields ...zap.Field) {
	printer.Error("["+tag+"] "+msg, fields...)
}

// AppStart
// 输出应用启动信息
func AppStart() {
	tag := "Boot"
	Info(tag, "============================================")
	Info(tag, "   ___                ____   _       _")
	Info(tag, "  / _ \\  _ __    ___ |  _ \\ (_) ___ | | __")
	Info(tag, " | | | || '_ \\  / _ \\| | | || |/ __|| |/ /")
	Info(tag, " | |_| || | | ||  __/| |_| || |\\__ \\|   <")
	Info(tag, "  \\___/ |_| |_| \\___||____/ |_||___/|_|\\_\\")
	Info(tag, "============================================")
	Info(tag, format.String("Welcome to OneDisk %s (%d)", definition.VersionName, definition.VersionCode))
	Info(tag, "Application is starting...")
}
