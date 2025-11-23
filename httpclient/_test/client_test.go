package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("创建客户端失败")
	}

	if client.retryConfig == nil {
		t.Error("默认配置未设置")
	}

	if client.retryConfig.MaxRetry != 3 {
		t.Errorf("默认最大重试次数不正确，期望: 3, 得到: %d", client.retryConfig.MaxRetry)
	}

	t.Log("New() 测试通过")
}

func TestClient_Get(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("请求方法不正确，期望: GET, 得到: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	}))
	defer server.Close()

	client := New()
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("GET请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("状态码不正确，期望: 200, 得到: %d", resp.StatusCode)
	}

	t.Log("Get() 测试通过")
}

func TestClient_Post(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("请求方法不正确，期望: POST, 得到: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post Success"))
	}))
	defer server.Close()

	client := New()
	resp, err := client.Post(server.URL, "text/plain", nil)
	if err != nil {
		t.Fatalf("POST请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("状态码不正确，期望: 200, 得到: %d", resp.StatusCode)
	}

	t.Log("Post() 测试通过")
}

func TestClient_Retry(t *testing.T) {
	attemptCount := 0

	// 创建测试服务器，前2次失败，第3次成功
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success on retry"))
	}))
	defer server.Close()

	// 创建客户端，配置重试
	config := &RetryConfig{
		MaxRetry:       3,
		RequestTimeout: 5 * time.Second,
		RetryDelay:     100 * time.Millisecond,
		EnableLog:      false,
	}
	client := NewWithConfig(nil, config)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("重试请求失败: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 3 {
		t.Errorf("重试次数不正确，期望: 3, 得到: %d", attemptCount)
	}

	t.Logf("重试测试通过，尝试次数: %d", attemptCount)
}

func TestClient_Timeout(t *testing.T) {
	// 创建慢响应的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建客户端，设置短超时时间
	config := &RetryConfig{
		MaxRetry:       1,
		RequestTimeout: 500 * time.Millisecond,
		EnableLog:      false,
	}
	client := NewWithConfig(nil, config)

	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("应该超时但没有返回错误")
	}

	t.Log("超时测试通过")
}

func TestClient_WithContext(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello Context"))
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	resp, err := client.DoWithContext(ctx, req)
	if err != nil {
		t.Fatalf("带上下文的请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("状态码不正确，期望: 200, 得到: %d", resp.StatusCode)
	}

	t.Log("WithContext() 测试通过")
}

func TestClient_ContextCancel(t *testing.T) {
	// 创建慢响应的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	_, err = client.DoWithContext(ctx, req)
	if err == nil {
		t.Error("上下文取消应该返回错误")
	}

	t.Log("上下文取消测试通过")
}

func TestClient_SetRetryConfig(t *testing.T) {
	client := New()

	newConfig := &RetryConfig{
		MaxRetry:       5,
		RequestTimeout: 60 * time.Second,
		EnableLog:      false,
	}

	client.SetRetryConfig(newConfig)

	if client.retryConfig.MaxRetry != 5 {
		t.Errorf("配置未正确更新，期望: 5, 得到: %d", client.retryConfig.MaxRetry)
	}

	t.Log("SetRetryConfig() 测试通过")
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config == nil {
		t.Fatal("默认配置为nil")
	}

	if config.MaxRetry != 3 {
		t.Errorf("默认MaxRetry不正确，期望: 3, 得到: %d", config.MaxRetry)
	}

	if config.RequestTimeout != 30*time.Second {
		t.Errorf("默认RequestTimeout不正确，期望: 30s, 得到: %v", config.RequestTimeout)
	}

	if config.RetryDelay != 1*time.Second {
		t.Errorf("默认RetryDelay不正确，期望: 1s, 得到: %v", config.RetryDelay)
	}

	if config.EnableLog != true {
		t.Error("默认EnableLog应该为true")
	}

	t.Log("DefaultRetryConfig() 测试通过")
}

func TestClient_MultipleRequests(t *testing.T) {
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request received"))
	}))
	defer server.Close()

	client := New()

	// 发送多个请求
	for i := 0; i < 5; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("第%d个请求失败: %v", i+1, err)
		}
		resp.Body.Close()
	}

	if requestCount != 5 {
		t.Errorf("请求次数不正确，期望: 5, 得到: %d", requestCount)
	}

	t.Log("多个请求测试通过")
}

func TestClient_DisableLog(t *testing.T) {
	config := &RetryConfig{
		MaxRetry:  1,
		EnableLog: false,
	}

	client := NewWithConfig(nil, config)

	if client.retryConfig.EnableLog {
		t.Error("日志应该被禁用")
	}

	t.Log("禁用日志测试通过")
}

// 性能测试
func BenchmarkClient_Get(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Benchmark"))
	}))
	defer server.Close()

	client := New()
	client.SetRetryConfig(&RetryConfig{
		MaxRetry:  1,
		EnableLog: false,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Get(server.URL)
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func BenchmarkClient_Post(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Benchmark"))
	}))
	defer server.Close()

	client := New()
	client.SetRetryConfig(&RetryConfig{
		MaxRetry:  1,
		EnableLog: false,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := client.Post(server.URL, "text/plain", nil)
		if resp != nil {
			resp.Body.Close()
		}
	}
}
