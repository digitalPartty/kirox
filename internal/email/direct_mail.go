package email

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DirectMailConfig struct {
	Name         string `json:"name"`
	BaseURL      string `json:"baseUrl"`
	RefreshToken string `json:"refreshToken"`
	ClientID     string `json:"clientId"`
	Email        string `json:"email"`
	Mailbox      string `json:"mailbox"`
}

type DirectMailClient struct {
	cfg    DirectMailConfig
	client *http.Client
}

type directMailResponse struct {
	Mails []struct {
		ID      interface{} `json:"id"`
		UID     interface{} `json:"uid"`
		Subject string      `json:"subject"`
		Text    string      `json:"text"`
		Body    string      `json:"body"`
		Raw     string      `json:"raw"`
	} `json:"mails"`
}

func NewDirectMailClient(cfg DirectMailConfig) *DirectMailClient {
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if cfg.Mailbox == "" {
		cfg.Mailbox = "INBOX"
	}
	return &DirectMailClient{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *DirectMailClient) fetchNewMails() (*directMailResponse, error) {
	params := url.Values{
		"refresh_token": {c.cfg.RefreshToken},
		"client_id":     {c.cfg.ClientID},
		"email":         {c.cfg.Email},
		"mailbox":       {c.cfg.Mailbox},
		"response_type": {"json"},
	}
	u := fmt.Sprintf("%s/api/mail-new?%s", c.cfg.BaseURL, params.Encode())

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败 (%d): %s", resp.StatusCode, string(body))
	}

	var result directMailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &result, nil
}

func (c *DirectMailClient) Test() error {
	_, err := c.fetchNewMails()
	return err
}

func (c *DirectMailClient) WaitForCode(timeoutSec, intervalSec int) (string, error) {
	maxAttempts := timeoutSec / intervalSec
	if maxAttempts <= 0 {
		maxAttempts = 1
	}
	processedIDs := make(map[string]bool)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		data, err := c.fetchNewMails()
		if err != nil {
			if attempt%5 == 0 {
				log.Printf("[DirectMail] 获取邮件失败: %v", err)
			}
			time.Sleep(time.Duration(intervalSec) * time.Second)
			continue
		}
		for _, mail := range data.Mails {
			mailID := fmt.Sprintf("%v", mail.ID)
			if mailID == "" || mailID == "<nil>" {
				mailID = fmt.Sprintf("%v", mail.UID)
			}
			if processedIDs[mailID] {
				continue
			}
			content := mail.Text
			if content == "" {
				content = mail.Body
			}
			if content == "" {
				content = mail.Raw
			}
			if code := ExtractVerificationCode(content); code != "" {
				log.Printf("[DirectMail] 找到验证码: %s", code)
				return code, nil
			}
			processedIDs[mailID] = true
		}
		if attempt < maxAttempts {
			time.Sleep(time.Duration(intervalSec) * time.Second)
		}
	}
	return "", fmt.Errorf("等待验证码超时 (%ds)", timeoutSec)
}

func (c *DirectMailClient) GetAddress() string {
	return c.cfg.Email
}
