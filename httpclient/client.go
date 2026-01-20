package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/efarxs/fastkit/httpclient/proxy"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetry        int           // 最大重试次数，默认3次
	RequestTimeout  time.Duration // 请求超时时间，默认30秒
	FallbackToLocal bool          // 失败后是否回退到本地请求，默认false
	RetryDelay      time.Duration // 重试延迟时间，默认1秒（每次重试会增加）
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetry:        3,
		RequestTimeout:  30 * time.Second,
		FallbackToLocal: false,
		RetryDelay:      1 * time.Second,
	}
}

// Client HTTP客户端封装
type Client struct {
	proxyManager *proxy.ProxyManager
	retryConfig  *RetryConfig
	logger       Logger // 日志记录器（可选）
}

// New 创建新的HTTP客户端（不使用代理）
func New() *Client {
	return &Client{
		retryConfig: DefaultRetryConfig(),
	}
}

// NewWithProxy 创建带代理的HTTP客户端
func NewWithProxy(proxyManager *proxy.ProxyManager) *Client {
	return &Client{
		proxyManager: proxyManager,
		retryConfig:  DefaultRetryConfig(),
		logger:       &NoOpLogger{}, // 默认不输出日志
	}
}

// NewWithConfig 创建带自定义配置的HTTP客户端
func NewWithConfig(proxyManager *proxy.ProxyManager, retryConfig *RetryConfig) *Client {
	if retryConfig == nil {
		retryConfig = DefaultRetryConfig()
	}
	return &Client{
		proxyManager: proxyManager,
		retryConfig:  retryConfig,
		logger:       &NoOpLogger{}, // 默认不输出日志
	}
}

// SetRetryConfig 设置重试配置
func (c *Client) SetRetryConfig(config *RetryConfig) {
	c.retryConfig = config
}

// SetProxyManager 设置代理管理器
func (c *Client) SetProxyManager(manager *proxy.ProxyManager) {
	c.proxyManager = manager
}

// SetLogger 设置日志记录器
func (c *Client) SetLogger(logger Logger) {
	if logger == nil {
		c.logger = &NoOpLogger{}
	} else {
		c.logger = logger
	}
}

// Do 执行HTTP请求（带重试和代理）
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.DoWithContext(context.Background(), req)
}

// DoWithContext 执行HTTP请求（带上下文）
func (c *Client) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	maxRetry := c.retryConfig.MaxRetry
	if maxRetry <= 0 {
		maxRetry = 1
	}

	var lastErr error
	useProxy := c.proxyManager != nil && c.proxyManager.IsEnabled()

	for attempt := 1; attempt <= maxRetry; attempt++ {
		if c.logger != nil {
			c.logger.Debug("执行HTTP请求",
				"method", req.Method,
				"url", req.URL.String(),
				"attempt", attempt,
				"use_proxy", useProxy)
		}

		// 创建HTTP客户端
		client, proxyURL, err := c.createClient(useProxy)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("创建HTTP客户端失败", "attempt", attempt, "error", err)
			}
			lastErr = err

			// 如果是最后一次尝试，且配置了回退到本地，则尝试不用代理
			if attempt == maxRetry && c.retryConfig.FallbackToLocal && useProxy {
				if c.logger != nil {
					c.logger.Info("达到最大重试次数，切换到本地请求")
				}
				useProxy = false
				maxRetry++ // 增加一次本地请求的机会
				continue
			}

			time.Sleep(c.retryConfig.RetryDelay * time.Duration(attempt))
			continue
		}

		// 如果使用代理，输出代理信息
		if useProxy && proxyURL != "" && c.logger != nil {
			c.logger.Info("使用代理IP发送请求",
				"proxy", proxyURL,
				"method", req.Method,
				"url", req.URL.String(),
				"attempt", attempt)
		}

		// 执行请求
		resp, err := c.doRequest(ctx, client, req)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("HTTP请求失败",
					"attempt", attempt,
					"use_proxy", useProxy,
					"error", err)
			}
			lastErr = err

			// 如果使用代理失败，且已经重试多次，考虑释放代理
			if useProxy && attempt >= 2 && c.proxyManager != nil {
				if c.logger != nil {
					c.logger.Warn("代理请求多次失败，释放当前代理IP以便下次获取新IP")
				}
				c.proxyManager.ReleaseProxy()
			}

			// 如果是最后一次尝试，且配置了回退到本地，则尝试不用代理
			if attempt == maxRetry && c.retryConfig.FallbackToLocal && useProxy {
				if c.logger != nil {
					c.logger.Info("达到最大重试次数，切换到本地请求")
				}
				useProxy = false
				maxRetry++ // 增加一次本地请求的机会
				continue
			}

			time.Sleep(c.retryConfig.RetryDelay * time.Duration(attempt))
			continue
		}

		// 成功返回
		if useProxy && proxyURL != "" && c.logger != nil {
			c.logger.Info("代理请求成功",
				"proxy", proxyURL,
				"status_code", resp.StatusCode,
				"attempt", attempt)
		}
		return resp, nil
	}

	return nil, fmt.Errorf("HTTP请求失败，已重试%d次: %w", maxRetry, lastErr)
}

// createClient 创建HTTP客户端，返回客户端和代理URL
func (c *Client) createClient(useProxy bool) (*http.Client, string, error) {
	timeout := c.retryConfig.RequestTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
	}

	proxyURLStr := ""
	if useProxy && c.proxyManager != nil {
		proxyURL, err := c.proxyManager.GetProxyURL()
		if err != nil {
			return nil, "", fmt.Errorf("获取代理URL失败: %w", err)
		}

		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, "", fmt.Errorf("解析代理URL失败: %w", err)
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURLParsed),
		}

		proxyURLStr = proxyURL
	}

	return client, proxyURLStr, nil
}

// createClientWithoutRedirect 创建不跟随重定向的HTTP客户端
func (c *Client) createClientWithoutRedirect(useProxy bool) (*http.Client, string, error) {
	client, proxyURLStr, err := c.createClient(useProxy)
	if err != nil {
		return nil, "", err
	}

	// 禁用自动重定向
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return client, proxyURLStr, nil
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	// 添加上下文
	reqWithContext := req.WithContext(ctx)

	// 执行请求
	resp, err := client.Do(reqWithContext)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

// Get 执行GET请求
func (c *Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post 执行POST请求
func (c *Client) Post(urlStr string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// DoWithoutRedirect 执行HTTP请求但不跟随重定向
func (c *Client) DoWithoutRedirect(req *http.Request) (*http.Response, error) {
	return c.DoWithoutRedirectContext(context.Background(), req)
}

// DoWithoutRedirectContext 执行HTTP请求但不跟随重定向（带上下文）
func (c *Client) DoWithoutRedirectContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	maxRetry := c.retryConfig.MaxRetry
	if maxRetry <= 0 {
		maxRetry = 1
	}

	var lastErr error
	useProxy := c.proxyManager != nil && c.proxyManager.IsEnabled()

	for attempt := 1; attempt <= maxRetry; attempt++ {
		if c.logger != nil {
			c.logger.Debug("执行HTTP请求(不跟随重定向)",
				"method", req.Method,
				"url", req.URL.String(),
				"attempt", attempt,
				"use_proxy", useProxy)
		}

		// 创建不跟随重定向的HTTP客户端
		client, proxyURL, err := c.createClientWithoutRedirect(useProxy)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("创建HTTP客户端失败", "attempt", attempt, "error", err)
			}
			lastErr = err

			// 如果是最后一次尝试，且配置了回退到本地，则尝试不用代理
			if attempt == maxRetry && c.retryConfig.FallbackToLocal && useProxy {
				if c.logger != nil {
					c.logger.Info("达到最大重试次数，切换到本地请求")
				}
				useProxy = false
				maxRetry++ // 增加一次本地请求的机会
				continue
			}

			time.Sleep(c.retryConfig.RetryDelay * time.Duration(attempt))
			continue
		}

		// 如果使用代理，输出代理信息
		if useProxy && proxyURL != "" && c.logger != nil {
			c.logger.Info("使用代理IP发送请求(不跟随重定向)",
				"proxy", proxyURL,
				"method", req.Method,
				"url", req.URL.String(),
				"attempt", attempt)
		}

		// 执行请求
		resp, err := c.doRequest(ctx, client, req)
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("HTTP请求失败",
					"attempt", attempt,
					"use_proxy", useProxy,
					"error", err)
			}
			lastErr = err

			// 如果使用代理失败，且已经重试多次，考虑释放代理
			if useProxy && attempt >= 2 && c.proxyManager != nil {
				if c.logger != nil {
					c.logger.Warn("代理请求多次失败，释放当前代理IP以便下次获取新IP")
				}
				c.proxyManager.ReleaseProxy()
			}

			// 如果是最后一次尝试，且配置了回退到本地，则尝试不用代理
			if attempt == maxRetry && c.retryConfig.FallbackToLocal && useProxy {
				if c.logger != nil {
					c.logger.Info("达到最大重试次数，切换到本地请求")
				}
				useProxy = false
				maxRetry++ // 增加一次本地请求的机会
				continue
			}

			time.Sleep(c.retryConfig.RetryDelay * time.Duration(attempt))
			continue
		}

		// 成功返回
		if useProxy && proxyURL != "" && c.logger != nil {
			c.logger.Info("代理请求成功(不跟随重定向)",
				"proxy", proxyURL,
				"status_code", resp.StatusCode,
				"attempt", attempt)
		}
		return resp, nil
	}

	return nil, fmt.Errorf("HTTP请求失败，已重试%d次: %w", maxRetry, lastErr)
}
