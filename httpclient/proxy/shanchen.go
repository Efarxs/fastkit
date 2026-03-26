package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// ShanChenProvider 闪臣代理提供商
type ShanChenProvider struct {
	ApiKey   string // API密钥
	Time     string // 时间参数
	Account  string // 代理账号（如果需要认证）
	Password string // 代理密码（如果需要认证）
	ApiUrl   string // 提取链接（如果有，优先使用这个）
}

// NewShanChenProvider 创建闪臣代理提供商
func NewShanChenProvider(apiKey, timeParam, account, password, apiUrl string) *ShanChenProvider {
	return &ShanChenProvider{
		ApiKey:   apiKey,
		Time:     timeParam,
		Account:  account,
		Password: password,
		ApiUrl:   apiUrl,
	}
}

// GetProxy 获取代理IP
func (p *ShanChenProvider) GetProxy() (*ProxyData, error) {
	var u string
	if p.ApiUrl != "" {
		u = p.ApiUrl
	} else {
		u = fmt.Sprintf("https://sch.shanchendaili.com/api.html?action=get_ip&key=%s&time=%s&count=1&type=json&only=0",
			p.ApiKey, p.Time)
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("获取闪臣代理IP失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取闪臣响应失败: %w", err)
	}

	res := gjson.ParseBytes(body)
	if res.Get("count").String() != "1" {
		log.Println("闪臣代理返回:", string(body))
		msg := res.Get("info").String()
		if msg == "" {
			msg = "未知错误"
		}
		return nil, errors.New(msg)
	}

	pT, err := strconv.Atoi(p.Time)
	if err != nil {
		pT = 5
	}

	return &ProxyData{
		IP:       res.Get("list.0.sever").String(),
		Port:     res.Get("list.0.port").String(),
		Account:  p.Account,
		Password: p.Password,
		ExpireAt: time.Now().Add(time.Duration(pT) * time.Minute),
	}, nil
}

// GetName 获取提供商名称
func (p *ShanChenProvider) GetName() string {
	return "闪臣(ShanChen)"
}
