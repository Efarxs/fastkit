package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// HskProvider 花生壳代理提供商
type HskProvider struct {
	ApiKey   string // API密钥
	Time     string // 时间参数
	Account  string // 代理账号（如果需要认证）
	Password string // 代理密码（如果需要认证）
}

// NewHskProvider 创建花生壳代理提供商
func NewHskProvider(apiKey, timeParam, account, password string) *HskProvider {
	return &HskProvider{
		ApiKey:   apiKey,
		Time:     timeParam,
		Account:  account,
		Password: password,
	}
}

// GetProxy 获取代理IP
func (p *HskProvider) GetProxy() (*ProxyData, error) {
	u := fmt.Sprintf("https://getip.huashengdaili.com/servers.php?session=%s&time=%s&count=1&type=json&only=1&pw=no&protocol=socket&separator=1&iptype=tunnel&format=city,time&dev=web",
		p.ApiKey, p.Time)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("获取花生壳代理IP失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取花生壳响应失败: %w", err)
	}

	// {"status":"0","count":1,"list":[{"sever":"113.14.131.59","port":2324,"net_type":1,"expire_time":"2025-11-06 22:13:05","province":"广西","city":"桂林市"}]}
	res := gjson.ParseBytes(body)
	if res.Get("status").String() != "0" {
		log.Println("花生壳代理返回:", string(body))
		msg := res.Get("msg").String()
		if msg == "" {
			msg = "未知错误"
		}
		return nil, errors.New(msg)
	}

	expireTime, _ := time.Parse("2006-01-02 15:04:05", res.Get("list.0.expire_time").String())
	
	return &ProxyData{
		IP:       res.Get("list.0.sever").String(),
		Port:     res.Get("list.0.port").String(),
		Account:  p.Account,
		Password: p.Password,
		ExpireAt: expireTime,
	}, nil
}

// GetName 获取提供商名称
func (p *HskProvider) GetName() string {
	return "花生壳(HSK)"
}
