package email

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DuckDuckGoConfig struct {
	Token string `json:"token"`
}

type DuckDuckGoClient struct {
	token  string
	client *http.Client
}

func NewDuckDuckGoClient(token string) *DuckDuckGoClient {
	return &DuckDuckGoClient{
		token:  strings.TrimSpace(token),
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *DuckDuckGoClient) CreateAlias() (string, error) {
	if c.token == "" {
		return "", fmt.Errorf("DuckDuckGo Token 未配置")
	}

	req, err := http.NewRequest("POST", "https://quack.duckduckgo.com/api/email/addresses", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("创建别名失败 (%d): %s", resp.StatusCode, string(body))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	var addr string
	for _, key := range []string{"address", "email", "alias", "duck_address", "private_address"} {
		if v, ok := data[key].(string); ok && v != "" {
			addr = v
			break
		}
	}
	if addr == "" {
		return "", fmt.Errorf("响应中未找到邮箱地址: %s", string(body))
	}
	if !strings.Contains(addr, "@") {
		addr = addr + "@duck.com"
	}
	return addr, nil
}

func (c *DuckDuckGoClient) Test() error {
	_, err := c.CreateAlias()
	return err
}
