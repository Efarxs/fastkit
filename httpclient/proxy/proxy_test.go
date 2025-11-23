package proxy

import (
	"testing"
	"time"
)

// MockProvider 模拟代理提供商
type MockProvider struct {
	ip       string
	port     string
	callCount int
}

func (m *MockProvider) GetProxy() (*ProxyData, error) {
	m.callCount++
	return &ProxyData{
		IP:       m.ip,
		Port:     m.port,
		Account:  "testuser",
		Password: "testpass",
		ExpireAt: time.Now().Add(5 * time.Minute),
	}, nil
}

func (m *MockProvider) GetName() string {
	return "MockProvider"
}

func TestNewProxyManager(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	if manager == nil {
		t.Fatal("创建代理管理器失败")
	}

	if !manager.IsEnabled() {
		t.Error("代理应该被启用")
	}

	t.Log("NewProxyManager() 测试通过")
}

func TestProxyManager_GetProxyURL(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	url, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("获取代理URL失败: %v", err)
	}

	expected := "http://testuser:testpass@1.2.3.4:8080"
	if url != expected {
		t.Errorf("代理URL不正确，期望: %s, 得到: %s", expected, url)
	}

	t.Logf("代理URL: %s", url)
	t.Log("GetProxyURL() 测试通过")
}

// MockProviderNoAuth 无认证的模拟代理提供商
type MockProviderNoAuth struct {
	ip   string
	port string
}

func (m *MockProviderNoAuth) GetProxy() (*ProxyData, error) {
	return &ProxyData{
		IP:       m.ip,
		Port:     m.port,
		Account:  "",
		Password: "",
		ExpireAt: time.Now().Add(5 * time.Minute),
	}, nil
}

func (m *MockProviderNoAuth) GetName() string {
	return "MockProviderNoAuth"
}

func TestProxyManager_GetProxyURL_NoAuth(t *testing.T) {
	provider := &MockProviderNoAuth{ip: "5.6.7.8", port: "9090"}

	manager := NewProxyManager(provider, true)
	url, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("获取代理URL失败: %v", err)
	}

	expected := "http://5.6.7.8:9090"
	if url != expected {
		t.Errorf("无认证代理URL不正确，期望: %s, 得到: %s", expected, url)
	}

	t.Log("无认证代理URL测试通过")
}

func TestProxyManager_Caching(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	// 第一次获取
	url1, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("第一次获取代理URL失败: %v", err)
	}

	// 第二次获取（应该使用缓存）
	url2, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("第二次获取代理URL失败: %v", err)
	}

	if url1 != url2 {
		t.Errorf("缓存的代理URL不一致")
	}

	if provider.callCount != 1 {
		t.Errorf("提供商调用次数不正确，期望: 1, 得到: %d", provider.callCount)
	}

	t.Log("代理缓存测试通过")
}

func TestProxyManager_SetEnabled(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	if !manager.IsEnabled() {
		t.Error("代理应该被启用")
	}

	manager.SetEnabled(false)

	if manager.IsEnabled() {
		t.Error("代理应该被禁用")
	}

	manager.SetEnabled(true)

	if !manager.IsEnabled() {
		t.Error("代理应该被重新启用")
	}

	t.Log("SetEnabled() 测试通过")
}

func TestProxyManager_SetProvider(t *testing.T) {
	provider1 := &MockProvider{ip: "1.2.3.4", port: "8080"}
	provider2 := &MockProvider{ip: "5.6.7.8", port: "9090"}

	manager := NewProxyManager(provider1, true)

	url1, _ := manager.GetProxyURL()

	// 切换提供商
	manager.SetProvider(provider2)

	url2, _ := manager.GetProxyURL()

	if url1 == url2 {
		t.Error("切换提供商后URL应该不同")
	}

	t.Logf("提供商1: %s", url1)
	t.Logf("提供商2: %s", url2)
	t.Log("SetProvider() 测试通过")
}

func TestProxyManager_ReleaseProxy(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	// 获取代理
	_, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("获取代理URL失败: %v", err)
	}

	// 释放代理
	manager.ReleaseProxy()

	// 再次获取（应该调用provider）
	_, err = manager.GetProxyURL()
	if err != nil {
		t.Fatalf("释放后获取代理URL失败: %v", err)
	}

	if provider.callCount != 2 {
		t.Errorf("提供商调用次数不正确，期望: 2, 得到: %d", provider.callCount)
	}

	t.Log("ReleaseProxy() 测试通过")
}

func TestProxyManager_GetCurrentProxy(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	// 初始应该没有代理
	currentProxy := manager.GetCurrentProxy()
	if currentProxy != nil {
		t.Error("初始应该没有当前代理")
	}

	// 获取代理
	_, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("获取代理URL失败: %v", err)
	}

	// 现在应该有代理
	currentProxy = manager.GetCurrentProxy()
	if currentProxy == nil {
		t.Error("应该有当前代理")
	}

	if currentProxy.IP != "1.2.3.4" {
		t.Errorf("当前代理IP不正确，期望: 1.2.3.4, 得到: %s", currentProxy.IP)
	}

	t.Log("GetCurrentProxy() 测试通过")
}

func TestProxyManager_DisabledProvider(t *testing.T) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, false) // 禁用

	if manager.IsEnabled() {
		t.Error("代理不应该被启用")
	}

	t.Log("禁用代理提供商测试通过")
}

func TestProxyManager_NilProvider(t *testing.T) {
	manager := NewProxyManager(nil, true)

	if manager.IsEnabled() {
		t.Error("没有提供商时不应该启用")
	}

	_, err := manager.GetProxyURL()
	if err == nil {
		t.Error("没有提供商时应该返回错误")
	}

	t.Log("nil提供商测试通过")
}

// MockExpiredProvider 返回已过期代理的提供商
type MockExpiredProvider struct {
	ip        string
	port      string
	callCount int
}

func (m *MockExpiredProvider) GetProxy() (*ProxyData, error) {
	m.callCount++
	return &ProxyData{
		IP:       m.ip,
		Port:     m.port,
		Account:  "",
		Password: "",
		ExpireAt: time.Now().Add(-1 * time.Minute), // 已过期
	}, nil
}

func (m *MockExpiredProvider) GetName() string {
	return "MockExpiredProvider"
}

func TestProxyManager_ExpiredProxy(t *testing.T) {
	// 创建一个返回已过期代理的提供商
	expiredProvider := &MockExpiredProvider{ip: "1.2.3.4", port: "8080"}

	manager := NewProxyManager(expiredProvider, true)

	// 第一次获取
	_, err := manager.GetProxyURL()
	if err != nil {
		t.Fatalf("第一次获取失败: %v", err)
	}

	// 等待一小段时间
	time.Sleep(100 * time.Millisecond)

	// 第二次获取（由于已过期，应该重新调用provider）
	_, err = manager.GetProxyURL()
	if err != nil {
		t.Fatalf("第二次获取失败: %v", err)
	}

	if expiredProvider.callCount < 2 {
		t.Errorf("过期代理应该重新获取，调用次数: %d", expiredProvider.callCount)
	}

	t.Log("过期代理测试通过")
}

// 性能测试
func BenchmarkProxyManager_GetProxyURL(b *testing.B) {
	provider := &MockProvider{ip: "1.2.3.4", port: "8080"}
	manager := NewProxyManager(provider, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GetProxyURL()
	}
}

func BenchmarkProxyManager_SetProvider(b *testing.B) {
	provider1 := &MockProvider{ip: "1.2.3.4", port: "8080"}
	provider2 := &MockProvider{ip: "5.6.7.8", port: "9090"}
	manager := NewProxyManager(provider1, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			manager.SetProvider(provider1)
		} else {
			manager.SetProvider(provider2)
		}
	}
}
