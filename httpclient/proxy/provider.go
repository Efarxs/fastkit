package proxy

import (
	"time"
)

// ProxyData 代理IP数据
type ProxyData struct {
	IP       string    // 代理IP地址
	Port     string    // 代理端口
	Account  string    // 代理账号（如果需要认证）
	Password string    // 代理密码（如果需要认证）
	ExpireAt time.Time // 过期时间
}

// ProxyProvider 代理提供商接口
type ProxyProvider interface {
	// GetProxy 获取代理IP
	GetProxy() (*ProxyData, error)

	// GetName 获取提供商名称
	GetName() string
}
