# 迁移指南

本文档帮助你从旧版本的httpclient迁移到新的通用版本。

## 主要变化

### 1. 移除了外部依赖

**旧版本**依赖：
- `wps-member-mod/config`
- `wps-member-mod/logger`
- `wps-member-mod/proxy`

**新版本**：
- 完全独立，无外部项目依赖
- 使用标准库的`log`包（可选）
- 所有配置通过参数传递

### 2. 代理配置方式改变

#### 旧版本
```go
// 旧版本依赖全局配置
client := httpclient.New()
// 配置在config文件中设置
```

#### 新版本
```go
// 方式1: 不使用代理
client := httpclient.New()

// 方式2: 使用代理
provider := proxy.NewHskProvider("api_key", "60", "user", "pass")
proxyManager := proxy.NewProxyManager(provider, true)
client := httpclient.NewWithProxy(proxyManager)

// 方式3: 自定义配置
retryConfig := &httpclient.RetryConfig{
    MaxRetry:        5,
    RequestTimeout:  60 * time.Second,
    FallbackToLocal: true,
    RetryDelay:      2 * time.Second,
    EnableLog:       true,
}
client := httpclient.NewWithConfig(proxyManager, retryConfig)
```

### 3. 代理提供商接口化

#### 旧版本
```go
// 旧版本使用固定的代理API
// 在proxy.go中硬编码
```

#### 新版本
```go
// 实现ProxyProvider接口即可添加新的代理商
type MyProvider struct {
    ApiKey string
}

func (p *MyProvider) GetProxy() (*proxy.ProxyData, error) {
    // 实现获取代理逻辑
}

func (p *MyProvider) GetName() string {
    return "MyProxy"
}

// 使用
myProvider := &MyProvider{ApiKey: "xxx"}
proxyManager := proxy.NewProxyManager(myProvider, true)
```

### 4. 重试配置改变

#### 旧版本
```go
// 旧版本在config文件中配置
type ProxyConfig struct {
    MaxRetry        int
    RequestTimeout  int
    FallbackToLocal bool
    // ...
}
```

#### 新版本
```go
// 新版本通过RetryConfig结构体配置
retryConfig := &httpclient.RetryConfig{
    MaxRetry:        3,
    RequestTimeout:  30 * time.Second,
    FallbackToLocal: false,
    RetryDelay:      1 * time.Second,
    EnableLog:       true,
}

client.SetRetryConfig(retryConfig)
```

## 迁移步骤

### 第一步：更新导入

```go
// 旧版本
import (
    "your-project/httpclient"
    "your-project/proxy"
)

// 新版本
import (
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/httpclient/proxy"
)
```

### 第二步：替换客户端创建方式

```go
// 旧版本
client := httpclient.New()

// 新版本（不使用代理）
client := httpclient.New()

// 新版本（使用代理）
provider := proxy.NewHskProvider("key", "60", "", "")
proxyManager := proxy.NewProxyManager(provider, true)
client := httpclient.NewWithProxy(proxyManager)
```

### 第三步：迁移配置

如果你之前使用配置文件，需要将配置转换为代码：

```go
// 旧版本 config.yaml
// proxy:
//   enabled: true
//   order_id: "xxx"
//   secret: "yyy"
//   max_retry: 3
//   request_timeout: 30

// 新版本代码
provider := proxy.NewHskProvider(orderID, "60", "", "")
proxyManager := proxy.NewProxyManager(provider, true)

retryConfig := &httpclient.RetryConfig{
    MaxRetry:       3,
    RequestTimeout: 30 * time.Second,
}

client := httpclient.NewWithConfig(proxyManager, retryConfig)
```

### 第四步：更新日志

```go
// 旧版本使用logger包
logger.HttpInfo("消息")

// 新版本
// 1. 使用内置日志（默认启用）
retryConfig := &httpclient.RetryConfig{
    EnableLog: true,
}

// 2. 禁用内置日志，使用自己的日志系统
retryConfig := &httpclient.RetryConfig{
    EnableLog: false,
}
// 然后在你的代码中添加日志
```

## 兼容性说明

### 完全兼容的功能
- HTTP请求的基本方法（Get, Post, Do等）
- 代理IP自动切换
- 重试机制
- 超时控制

### 需要调整的功能
- 配置方式（从配置文件改为代码配置）
- 日志输出（从zap改为标准log）
- 代理管理器获取方式（从单例改为手动创建）

### 新增功能
- 支持多种代理提供商（花生壳、巨量、闪臣等）
- 可自定义代理提供商
- 动态切换代理提供商
- 更灵活的重试配置
- 可选的日志输出

## 常见问题

### Q1: 如何保持与旧版本相同的行为？

```go
// 创建类似旧版本的客户端
provider := proxy.NewHskProvider(orderID, "60", "", "")
proxyManager := proxy.NewProxyManager(provider, true)

retryConfig := &httpclient.RetryConfig{
    MaxRetry:        3,
    RequestTimeout:  30 * time.Second,
    FallbackToLocal: true,
    RetryDelay:      1 * time.Second,
    EnableLog:       true,
}

client := httpclient.NewWithConfig(proxyManager, retryConfig)
```

### Q2: 如何在运行时切换代理商？

```go
// 初始化
proxyManager := proxy.NewProxyManager(hskProvider, true)
client := httpclient.NewWithProxy(proxyManager)

// 运行时切换
proxyManager.SetProvider(jlProvider)
```

### Q3: 如何禁用日志？

```go
retryConfig := &httpclient.RetryConfig{
    EnableLog: false,
}
client.SetRetryConfig(retryConfig)
```

### Q4: 如何添加新的代理商？

参考`hsk.go`、`jl.go`、`shanchen.go`的实现，创建一个实现了`ProxyProvider`接口的结构体即可。

## 建议

1. **逐步迁移**：先在新功能中使用新版本，再逐步迁移旧代码
2. **配置管理**：可以创建一个配置管理函数来统一创建客户端
3. **错误处理**：新版本的错误信息更详细，建议检查错误处理逻辑
4. **测试**：迁移后充分测试，特别是代理切换和重试逻辑

## 示例：完整的迁移案例

### 旧版本代码
```go
package main

import "your-project/httpclient"

func main() {
    client := httpclient.New()
    resp, err := client.Get("https://api.example.com/data")
    // ...
}
```

### 新版本代码
```go
package main

import (
    "time"
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/httpclient/proxy"
)

func main() {
    // 创建代理提供商
    provider := proxy.NewHskProvider(
        "your_api_key",
        "60",
        "",
        "",
    )
    
    // 创建代理管理器
    proxyManager := proxy.NewProxyManager(provider, true)
    
    // 配置重试
    retryConfig := &httpclient.RetryConfig{
        MaxRetry:        3,
        RequestTimeout:  30 * time.Second,
        FallbackToLocal: true,
        RetryDelay:      1 * time.Second,
        EnableLog:       true,
    }
    
    // 创建客户端
    client := httpclient.NewWithConfig(proxyManager, retryConfig)
    
    // 发送请求
    resp, err := client.Get("https://api.example.com/data")
    // ...
}
```
