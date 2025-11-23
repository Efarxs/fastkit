# HttpClient - 通用HTTP客户端库

一个支持代理IP、重试机制、可配置的Go HTTP客户端库。

## 特性

- ✅ 支持自定义重试配置（次数、超时、延迟等）
- 🌐 支持多种代理提供商（花生壳、巨量、闪臣等）
- 🔧 易于扩展，可自定义代理提供商
- 🔄 自动重试和代理IP切换
- 📝 可选的日志输出
- 🛡️ 代理失败自动回退到本地请求

## 安装

```bash
go get github.com/efarxs/fastkit/httpclient
```

## 快速开始

### 1. 不使用代理的简单HTTP请求

```go
package main

import (
    "fmt"
    "io"
    "github.com/efarxs/fastkit/httpclient"
)

func main() {
    // 创建客户端
    client := httpclient.New()
    
    // 发送GET请求
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 2. 使用代理IP（花生壳示例）

```go
package main

import (
    "fmt"
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/httpclient/proxy"
)

func main() {
    // 创建花生壳代理提供商
    hskProvider := proxy.NewHskProvider(
        "your_api_key",    // API密钥
        "60",              // 时间参数（分钟）
        "username",        // 账号（如果需要）
        "password",        // 密码（如果需要）
    )
    
    // 创建代理管理器
    proxyManager := proxy.NewProxyManager(hskProvider, true)
    
    // 创建带代理的HTTP客户端
    client := httpclient.NewWithProxy(proxyManager)
    
    // 发送请求
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    fmt.Println("状态码:", resp.StatusCode)
}
```

### 3. 使用巨量代理

```go
// 创建巨量代理提供商
jlProvider := proxy.NewJlProvider(
    "https://api.jldaili.com/getip?...", // 巨量API完整URL
    "username",                           // 账号
    "password",                           // 密码
)

proxyManager := proxy.NewProxyManager(jlProvider, true)
client := httpclient.NewWithProxy(proxyManager)
```

### 4. 使用闪臣代理

```go
// 创建闪臣代理提供商
shanChenProvider := proxy.NewShanChenProvider(
    "your_api_key",  // API密钥
    "5",             // 时间参数（分钟）
    "username",      // 账号
    "password",      // 密码
)

proxyManager := proxy.NewProxyManager(shanChenProvider, true)
client := httpclient.NewWithProxy(proxyManager)
```

### 5. 自定义重试配置

```go
package main

import (
    "time"
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/httpclient/proxy"
)

func main() {
    // 自定义重试配置
    retryConfig := &httpclient.RetryConfig{
        MaxRetry:        5,               // 最大重试5次
        RequestTimeout:  60 * time.Second, // 请求超时60秒
        FallbackToLocal: true,             // 代理失败后回退到本地
        RetryDelay:      2 * time.Second,  // 重试延迟2秒
    }
    
    // 创建代理（可选）
    provider := proxy.NewHskProvider("api_key", "60", "", "")
    proxyManager := proxy.NewProxyManager(provider, true)
    
    // 创建带自定义配置的客户端
    client := httpclient.NewWithConfig(proxyManager, retryConfig)
    
    // 启用日志（可选）
    client.SetLogger(httpclient.NewDefaultLogger())
    
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### 6. 使用自定义日志（集成项目logger）

```go
package main

import (
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/logger"
)

func main() {
    // 初始化项目logger
    logger.Init()
    
    // 创建客户端
    client := httpclient.New()
    
    // 使用项目的zap logger
    client.SetLogger(httpclient.NewZapLogger(logger.HttpLogger))
    
    // 现在所有HTTP请求日志会输出到http.log
    resp, _ := client.Get("https://api.example.com/data")
    defer resp.Body.Close()
}
```

### 7. 动态切换代理提供商

```go
package main

import (
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/httpclient/proxy"
)

func main() {
    // 初始使用花生壳
    hskProvider := proxy.NewHskProvider("key1", "60", "", "")
    proxyManager := proxy.NewProxyManager(hskProvider, true)
    client := httpclient.NewWithProxy(proxyManager)
    
    // 发送请求
    client.Get("https://api.example.com/data1")
    
    // 切换到巨量代理
    jlProvider := proxy.NewJlProvider("https://api.jldaili.com/...", "", "")
    proxyManager.SetProvider(jlProvider)
    
    // 再次发送请求，使用新的代理
    client.Get("https://api.example.com/data2")
    
    // 禁用代理
    proxyManager.SetEnabled(false)
    client.Get("https://api.example.com/data3")
}
```

## 自定义代理提供商

你可以轻松实现自己的代理提供商：

```go
package proxy

import (
    "fmt"
    "time"
    "github.com/efarxs/fastkit/httpclient/proxy"
)

// MyCustomProvider 自定义代理提供商
type MyCustomProvider struct {
    ApiKey string
}

// GetProxy 实现ProxyProvider接口
func (p *MyCustomProvider) GetProxy() (*proxy.ProxyData, error) {
    // 调用你的代理API获取IP
    // ...
    
    return &proxy.ProxyData{
        IP:       "1.2.3.4",
        Port:     "8080",
        Account:  "user",
        Password: "pass",
        ExpireAt: time.Now().Add(5 * time.Minute),
    }, nil
}

// GetName 实现ProxyProvider接口
func (p *MyCustomProvider) GetName() string {
    return "MyCustomProxy"
}

// 使用自定义提供商
func main() {
    myProvider := &MyCustomProvider{ApiKey: "xxx"}
    proxyManager := proxy.NewProxyManager(myProvider, true)
    client := httpclient.NewWithProxy(proxyManager)
    
    client.Get("https://api.example.com/data")
}
```

## API文档

### RetryConfig 结构体

```go
type RetryConfig struct {
    MaxRetry        int           // 最大重试次数，默认3次
    RequestTimeout  time.Duration // 请求超时时间，默认30秒
    FallbackToLocal bool          // 失败后是否回退到本地请求，默认false
    RetryDelay      time.Duration // 重试延迟时间，默认1秒
}
```

### Logger 接口

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})  // 调试日志
    Info(msg string, keysAndValues ...interface{})   // 信息日志
    Warn(msg string, keysAndValues ...interface{})   // 警告日志
    Error(msg string, keysAndValues ...interface{})  // 错误日志
}
```

**内置日志实现：**
- `NewDefaultLogger()` - 使用标准库log的默认实现
- `NewZapLogger(logger)` - 基于zap的适配器，可集成项目logger
- `NoOpLogger` - 空实现，不输出任何日志

### ProxyProvider 接口

```go
type ProxyProvider interface {
    GetProxy() (*ProxyData, error) // 获取代理IP
    GetName() string                // 获取提供商名称
}
```

### Client 方法

- `New()` - 创建不带代理的客户端
- `NewWithProxy(proxyManager)` - 创建带代理的客户端
- `NewWithConfig(proxyManager, retryConfig)` - 创建带自定义配置的客户端
- `Get(url)` - 发送GET请求
- `Post(url, contentType, body)` - 发送POST请求
- `Do(req)` - 执行自定义请求
- `DoWithContext(ctx, req)` - 执行带上下文的请求
- `SetRetryConfig(config)` - 设置重试配置
- `SetProxyManager(manager)` - 设置代理管理器
- `SetLogger(logger)` - 设置日志记录器（传nil禁用日志）

### ProxyManager 方法

- `NewProxyManager(provider, enabled)` - 创建代理管理器
- `SetProvider(provider)` - 设置代理提供商
- `SetEnabled(enabled)` - 设置是否启用代理
- `IsEnabled()` - 是否启用代理
- `ReleaseProxy()` - 释放当前代理IP
- `GetCurrentProxy()` - 获取当前代理信息

## 注意事项

1. 代理账号和密码根据不同代理商要求填写，有些代理商不需要
2. 代理IP会自动管理过期时间，过期后自动获取新IP
3. 请求失败时会自动重试，并在重试多次后释放代理IP
4. 设置`FallbackToLocal`为`true`可在代理失败时自动切换到本地请求
5. 建议在生产环境中适当设置超时时间和重试次数

## License

MIT
