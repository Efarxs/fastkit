package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// JlProvider 巨量代理提供商
type JlProvider struct {
	ApiUrl   string // API URL（巨量直接把完整URL作为ApiUrl）
	Account  string // 代理账号（如果需要认证）
	Password string // 代理密码（如果需要认证）
}

// NewJlProvider 创建巨量代理提供商
func NewJlProvider(apiUrl, account, password string) *JlProvider {
	return &JlProvider{
		ApiUrl:   apiUrl,
		Account:  account,
		Password: password,
	}
}

// GetProxy 获取代理IP
func (p *JlProvider) GetProxy() (*ProxyData, error) {
	resp, err := http.Get(p.ApiUrl)
	if err != nil {
		return nil, fmt.Errorf("获取巨量代理IP失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取巨量响应失败: %w", err)
	}

	// {"code":200,"msg":"成功","data":{"count":1,"filter_count":1,"surplus_quantity":0,"proxy_list":["112.85.129.243:24333"]}}
	res := gjson.ParseBytes(body)
	if res.Get("code").String() != "200" {
		log.Println("巨量代理返回:", string(body))
		msg := res.Get("msg").String()
		if msg == "" {
			msg = "未知错误"
		}
		return nil, errors.New(msg)
	}

	ipPort := res.Get("data.proxy_list.0").String()
	if ipPort == "" {
		return nil, errors.New("未获取到代理IP")
	}

	parts := strings.Split(ipPort, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("代理IP格式错误: %s", ipPort)
	}

	return &ProxyData{
		IP:       parts[0],
		Port:     parts[1],
		Account:  p.Account,
		Password: p.Password,
		ExpireAt: time.Now().Add(5 * time.Minute), // 巨量默认5分钟过期
	}, nil
}

// GetName 获取提供商名称
func (p *JlProvider) GetName() string {
	return "巨量(JL)"
}
