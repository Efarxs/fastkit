# HttpClient 日志使用示例

## 方式1：默认不输出日志（推荐生产环境）

```go
client := httpclient.New()
// 默认不输出任何日志
resp, _ := client.Get("https://api.example.com")
```

## 方式2：使用标准库log

```go
client := httpclient.New()
client.SetLogger(httpclient.NewDefaultLogger())

// 日志会输出到标准输出
resp, _ := client.Get("https://api.example.com")
```

输出示例：
```
[HTTPClient] 2025/11/23 09:30:00 [DEBUG] 执行HTTP请求 [method GET url https://api.example.com attempt 1 use_proxy false]
[HTTPClient] 2025/11/23 09:30:00 [INFO] 代理请求成功 [proxy http://1.2.3.4:8080 status_code 200 attempt 1]
```

## 方式3：集成项目的zap logger（推荐）

```go
package main

import (
    "github.com/efarxs/fastkit/httpclient"
    "github.com/efarxs/fastkit/logger"
)

func main() {
    // 初始化项目logger
    logger.Init()
    
    // 创建HTTP客户端
    client := httpclient.New()
    
    // 使用项目的HTTP专用logger
    client.SetLogger(httpclient.NewZapLogger(logger.HttpLogger))
    
    // 所有HTTP请求日志会输出到 logs/http.log
    resp, _ := client.Get("https://api.example.com")
    defer resp.Body.Close()
}
```

日志输出（写入logs/http.log）：
```json
{"level":"debug","ts":1700712600,"msg":"执行HTTP请求","method":"GET","url":"https://api.example.com","attempt":1,"use_proxy":false}
{"level":"info","ts":1700712601,"msg":"代理请求成功","proxy":"http://1.2.3.4:8080","status_code":200,"attempt":1}
```

## 方式4：实现自定义Logger

```go
package main

import "github.com/efarxs/fastkit/httpclient"

// MyCustomLogger 自定义日志实现
type MyCustomLogger struct {
    // 你的日志字段
}

func (l *MyCustomLogger) Debug(msg string, keysAndValues ...interface{}) {
    // 实现Debug日志
}

func (l *MyCustomLogger) Info(msg string, keysAndValues ...interface{}) {
    // 实现Info日志
}

func (l *MyCustomLogger) Warn(msg string, keysAndValues ...interface{}) {
    // 实现Warn日志
}

func (l *MyCustomLogger) Error(msg string, keysAndValues ...interface{}) {
    // 实现Error日志
}

func main() {
    client := httpclient.New()
    client.SetLogger(&MyCustomLogger{})
    
    // 使用你的自定义logger
    resp, _ := client.Get("https://api.example.com")
    defer resp.Body.Close()
}
```

## 方式5：动态切换日志

```go
client := httpclient.New()

// 开发环境：启用详细日志
if isDevelopment {
    client.SetLogger(httpclient.NewDefaultLogger())
}

// 生产环境：使用项目logger
if isProduction {
    client.SetLogger(httpclient.NewZapLogger(logger.HttpLogger))
}

// 测试环境：禁用日志
if isTesting {
    client.SetLogger(nil)  // 或者不设置，默认就是NoOpLogger
}
```

## 日志级别说明

| 级别 | 使用场景 |
|-----|---------|
| Debug | 详细的请求信息（方法、URL、尝试次数等） |
| Info | 重要操作（使用代理、请求成功、回退到本地等） |
| Warn | 警告信息（创建客户端失败、请求失败、多次重试等） |
| Error | 错误信息（当前未使用，可扩展） |

## 最佳实践

### 1. 生产环境
```go
// 使用项目的logger，输出到文件
client.SetLogger(httpclient.NewZapLogger(logger.HttpLogger))
```

### 2. 开发环境
```go
// 使用默认logger，输出到控制台
client.SetLogger(httpclient.NewDefaultLogger())
```

### 3. 单元测试
```go
// 不输出日志，保持测试输出简洁
client := httpclient.New()  // 默认NoOpLogger
```

### 4. 性能敏感场景
```go
// 禁用日志以获得最佳性能
client.SetLogger(nil)
```
