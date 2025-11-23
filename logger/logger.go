package logger

import (
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var HttpLogger *zap.Logger // HTTP专用日志

// Init 初始化日志系统
func Init() {
	// 确保logs目录存在
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(err)
	}

	// 文件输出配置（使用 lumberjack 实现日志切割）
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logsDir, "app.log"), // 日志文件路径
		MaxSize:    2,                                 // 每个日志文件最大尺寸（MB）
		MaxBackups: 30,                                // 保留旧文件的最大个数
		MaxAge:     30,                                // 保留旧文件的最大天数
		Compress:   true,                              // 是否压缩/归档旧文件
		LocalTime:  true,                              // 使用本地时间
	}

	// 控制台输出配置
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色输出
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 文件输出配置
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 设置日志级别
	level := zapcore.InfoLevel

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),            // 控制台输出
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapcore.DebugLevel), // 文件输出（记录所有级别）
	)

	// 创建 logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 创建HTTP专用logger
	initHttpLogger()
}

// initHttpLogger 初始化HTTP专用日志
func initHttpLogger() {
	// 确保logs目录存在
	logsDir := "logs"

	// HTTP日志文件输出配置
	httpFileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logsDir, "http.log"), // HTTP日志文件路径
		MaxSize:    2,                                  // 每个日志文件最大尺寸（MB）
		MaxBackups: 10,                                 // 保留旧文件的最大个数
		MaxAge:     7,                                  // 保留旧文件的最大天数
		Compress:   true,                               // 是否压缩/归档旧文件
		LocalTime:  true,                               // 使用本地时间
	}

	// HTTP日志文件输出配置
	httpFileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// 创建HTTP日志核心（只输出到文件，不输出到控制台）
	httpCore := zapcore.NewCore(httpFileEncoder, zapcore.AddSync(httpFileWriter), zapcore.DebugLevel)

	// 创建HTTP logger
	HttpLogger = zap.New(httpCore, zap.AddCaller(), zap.AddCallerSkip(1))
}

// Sync 刷新日志缓冲区
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
	if HttpLogger != nil {
		_ = HttpLogger.Sync()
	}
}

// Info 记录Info级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug 记录Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn 记录Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 记录Error级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 记录Fatal级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Infof 格式化Info日志
func Infof(template string, args ...interface{}) {
	Logger.Sugar().Infof(template, args...)
}

// Debugf 格式化Debug日志
func Debugf(template string, args ...interface{}) {
	Logger.Sugar().Debugf(template, args...)
}

// Warnf 格式化Warn日志
func Warnf(template string, args ...interface{}) {
	Logger.Sugar().Warnf(template, args...)
}

// Errorf 格式化Error日志
func Errorf(template string, args ...interface{}) {
	Logger.Sugar().Errorf(template, args...)
}

// Fatalf 格式化Fatal日志并退出程序
func Fatalf(template string, args ...interface{}) {
	Logger.Sugar().Fatalf(template, args...)
}

// ==================== HTTP专用日志方法 ====================

// HttpInfo 记录HTTP Info级别日志
func HttpInfo(msg string, fields ...zap.Field) {
	HttpLogger.Info(msg, fields...)
}

// HttpDebug 记录HTTP Debug级别日志
func HttpDebug(msg string, fields ...zap.Field) {
	HttpLogger.Debug(msg, fields...)
}

// HttpWarn 记录HTTP Warn级别日志
func HttpWarn(msg string, fields ...zap.Field) {
	HttpLogger.Warn(msg, fields...)
}

// HttpError 记录HTTP Error级别日志
func HttpError(msg string, fields ...zap.Field) {
	HttpLogger.Error(msg, fields...)
}

// ==================== Panic恢复功能 ====================

// RecoverPanic 捕获panic并记录到日志
func RecoverPanic(ctx string) {
	if r := recover(); r != nil {
		// 获取堆栈信息
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]

		Error("捕获到异常panic",
			zap.String("context", ctx),
			zap.Any("panic", r),
			zap.String("stack", string(stack)))
	}
}
