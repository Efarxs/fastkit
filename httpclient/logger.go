package httpclient

import (
	"log"
	"os"
)

// Logger HTTP客户端日志接口
type Logger interface {
	// Debug 调试日志
	Debug(msg string, keysAndValues ...interface{})
	// Info 信息日志
	Info(msg string, keysAndValues ...interface{})
	// Warn 警告日志
	Warn(msg string, keysAndValues ...interface{})
	// Error 错误日志
	Error(msg string, keysAndValues ...interface{})
}

// DefaultLogger 默认日志实现（使用标准库log）
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志记录器
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "[HTTPClient] ", log.LstdFlags),
	}
}

func (l *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *DefaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Printf("[WARN] %s %v", msg, keysAndValues)
}

func (l *DefaultLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Printf("[ERROR] %s %v", msg, keysAndValues)
}

// NoOpLogger 空日志实现（不输出任何日志）
type NoOpLogger struct{}

func (l *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *NoOpLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *NoOpLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {}
