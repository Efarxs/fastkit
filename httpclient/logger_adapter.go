package httpclient

import (
	"go.uber.org/zap"
)

// ZapLoggerAdapter zap日志适配器
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

// NewZapLogger 创建基于zap的日志记录器
func NewZapLogger(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{
		logger: logger,
	}
}

func (l *ZapLoggerAdapter) Debug(msg string, keysAndValues ...interface{}) {
	fields := l.convertToZapFields(keysAndValues...)
	l.logger.Debug(msg, fields...)
}

func (l *ZapLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	fields := l.convertToZapFields(keysAndValues...)
	l.logger.Info(msg, fields...)
}

func (l *ZapLoggerAdapter) Warn(msg string, keysAndValues ...interface{}) {
	fields := l.convertToZapFields(keysAndValues...)
	l.logger.Warn(msg, fields...)
}

func (l *ZapLoggerAdapter) Error(msg string, keysAndValues ...interface{}) {
	fields := l.convertToZapFields(keysAndValues...)
	l.logger.Error(msg, fields...)
}

// convertToZapFields 转换键值对为zap字段
func (l *ZapLoggerAdapter) convertToZapFields(keysAndValues ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			break
		}
		
		key, ok := keysAndValues[i].(string)
		if !ok {
			continue
		}
		
		value := keysAndValues[i+1]
		fields = append(fields, zap.Any(key, value))
	}
	
	return fields
}
