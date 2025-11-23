package proxy

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ProxyManager 代理管理器
type ProxyManager struct {
	provider     ProxyProvider // 代理提供商
	currentProxy *ProxyData    // 当前代理IP
	enabled      bool          // 是否启用代理
	mu           sync.RWMutex  // 读写锁
}

// NewProxyManager 创建代理管理器
func NewProxyManager(provider ProxyProvider, enabled bool) *ProxyManager {
	return &ProxyManager{
		provider: provider,
		enabled:  enabled,
	}
}

// SetProvider 设置代理提供商
func (pm *ProxyManager) SetProvider(provider ProxyProvider) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.provider = provider
	pm.currentProxy = nil // 清空当前代理
}

// SetEnabled 设置是否启用代理
func (pm *ProxyManager) SetEnabled(enabled bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = enabled
}

// IsEnabled 是否启用代理
func (pm *ProxyManager) IsEnabled() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.enabled && pm.provider != nil
}

// GetProxyURL 获取代理URL（格式：http://ip:port 或 http://user:pass@ip:port）
func (pm *ProxyManager) GetProxyURL() (string, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.provider == nil {
		return "", fmt.Errorf("代理提供商未设置")
	}

	// 如果当前代理还有效，直接返回
	if pm.currentProxy != nil && time.Now().Before(pm.currentProxy.ExpireAt) {
		return pm.formatProxyURL(pm.currentProxy), nil
	}

	// 获取新的代理IP
	proxyData, err := pm.provider.GetProxy()
	if err != nil {
		return "", fmt.Errorf("获取代理IP失败: %w", err)
	}

	pm.currentProxy = proxyData
	
	log.Printf("获取到新的代理IP: %s:%s (来源: %s, 过期时间: %s)",
		proxyData.IP, proxyData.Port, pm.provider.GetName(), proxyData.ExpireAt.Format("2006-01-02 15:04:05"))

	return pm.formatProxyURL(proxyData), nil
}

// formatProxyURL 格式化代理URL
func (pm *ProxyManager) formatProxyURL(data *ProxyData) string {
	if data.Account != "" && data.Password != "" {
		return fmt.Sprintf("http://%s:%s@%s:%s", data.Account, data.Password, data.IP, data.Port)
	}
	return fmt.Sprintf("http://%s:%s", data.IP, data.Port)
}

// ReleaseProxy 释放当前代理
func (pm *ProxyManager) ReleaseProxy() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.currentProxy = nil
	log.Println("已释放当前代理IP")
}

// GetCurrentProxy 获取当前代理信息
func (pm *ProxyManager) GetCurrentProxy() *ProxyData {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.currentProxy
}
