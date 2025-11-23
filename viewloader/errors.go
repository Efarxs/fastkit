package viewloader

import (
	"errors"
	"fmt"
)

var (
	// ErrEmptyConfig 配置为空错误
	ErrEmptyConfig = errors.New("配置不能为空: TemplateGroups 和 StandalonePaths 至少需要配置一项")

	// ErrLayoutNotFound 布局模板文件不存在错误
	ErrLayoutNotFound = errors.New("布局模板文件不存在")

	// ErrViewNotFound 视图文件不存在错误
	ErrViewNotFound = errors.New("视图文件不存在")
)

// ConfigError 配置错误
type ConfigError struct {
	Field   string // 错误字段
	Message string // 错误消息
}

// Error 实现 error 接口
func (e *ConfigError) Error() string {
	return fmt.Sprintf("配置错误 [%s]: %s", e.Field, e.Message)
}

// LoadError 加载错误
type LoadError struct {
	Path    string // 文件路径
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *LoadError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("加载失败 [%s]: %s - %v", e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("加载失败 [%s]: %s", e.Path, e.Message)
}

// Unwrap 支持 errors.Unwrap
func (e *LoadError) Unwrap() error {
	return e.Err
}
