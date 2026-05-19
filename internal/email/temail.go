package email

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type TEmailConfig struct {
	Name          string `json:"name"`
	BaseURL       string `json:"baseUrl"`
	Email         string `json:"email"`
	JWT           string `json:"jwt,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
}

type TEmailClient struct {
	cfg    TEmailConfig
	client *http.Client
}

type temailMail struct {
	ID  int64  `json:"id"`
	Raw string `json:"raw"`
}

func NewTEmailClient(cfg TEmailConfig) *TEmailClient {
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	return &TEmailClient{
		cfg:    cfg,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *TEmailClient) request(method, urlStr string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

// fetchJWT 通过 Admin API 获取 JWT
func (c *TEmailClient) fetchJWT() (string, error) {
	if c.cfg.AdminPassword == "" {
		return "", fmt.Errorf("Admin 密码未设置")
	}
	emailName := c.cfg.Email
	if i := strings.Index(emailName, "@"); i > 0 {
		emailName = emailName[:i]
	}

	listURL := fmt.Sprintf("%s/admin/address?limit=100&offset=0&query=%s", c.cfg.BaseURL, url.QueryEscape(emailName))
	body, status, err := c.request("GET", listURL, map[string]string{"x-admin-auth": c.cfg.AdminPassword})
	if err != nil {
		return "", fmt.Errorf("查询邮箱失败: %w", err)
	}
	if status != 200 {
		return "", fmt.Errorf("查询邮箱失败 (%d): %s", status, string(body))
	}

	var listResp struct {
		Results []struct {
			ID int64 `json:"id"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &listResp); err != nil {
		return "", fmt.Errorf("解析邮箱列表失败: %w", err)
	}
	if len(listResp.Results) == 0 {
		return "", fmt.Errorf("找不到邮箱: %s", c.cfg.Email)
	}

	addrID := listResp.Results[0].ID
	jwtURL := fmt.Sprintf("%s/admin/show_password/%d", c.cfg.BaseURL, addrID)
	body, status, err = c.request("GET", jwtURL, map[string]string{"x-admin-auth": c.cfg.AdminPassword})
	if err != nil {
		return "", fmt.Errorf("获取 JWT 失败: %w", err)
	}
	if status != 200 {
		return "", fmt.Errorf("获取 JWT 失败 (%d): %s", status, string(body))
	}

	var jwtResp struct {
		JWT string `json:"jwt"`
	}
	if err := json.Unmarshal(body, &jwtResp); err != nil {
		return "", fmt.Errorf("解析 JWT 响应失败: %w", err)
	}
	if jwtResp.JWT == "" {
		return "", fmt.Errorf("响应中无 JWT 字段")
	}
	c.cfg.JWT = jwtResp.JWT
	return jwtResp.JWT, nil
}

func (c *TEmailClient) ensureJWT() error {
	if c.cfg.JWT != "" {
		return nil
	}
	if c.cfg.AdminPassword != "" {
		_, err := c.fetchJWT()
		return err
	}
	return fmt.Errorf("未配置 JWT 或 Admin 密码")
}

func (c *TEmailClient) fetchMails(limit, offset int) ([]temailMail, error) {
	if c.cfg.AdminPassword != "" {
		return c.fetchMailsViaAdmin(limit, offset)
	}
	if err := c.ensureJWT(); err != nil {
		return nil, err
	}
	return c.fetchMailsViaJWT(limit, offset)
}

func (c *TEmailClient) fetchMailsViaAdmin(limit, offset int) ([]temailMail, error) {
	u := fmt.Sprintf("%s/admin/mails?limit=%d&offset=%d", c.cfg.BaseURL, limit, offset)
	if c.cfg.Email != "" {
		u += "&address=" + url.QueryEscape(c.cfg.Email)
	}
	body, status, err := c.request("GET", u, map[string]string{"x-admin-auth": c.cfg.AdminPassword})
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("获取邮件失败 (%d)", status)
	}
	var resp struct {
		Results []temailMail `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析邮件列表失败: %w", err)
	}
	return resp.Results, nil
}

func (c *TEmailClient) fetchMailsViaJWT(limit, offset int) ([]temailMail, error) {
	u := fmt.Sprintf("%s/api/mails?limit=%d&offset=%d", c.cfg.BaseURL, limit, offset)
	body, status, err := c.request("GET", u, map[string]string{"Authorization": "Bearer " + c.cfg.JWT})
	if err != nil || status != 200 {
		return c.fetchMailsViaJWTParam(limit, offset)
	}
	var resp struct {
		Results []temailMail `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析邮件列表失败: %w", err)
	}
	if len(resp.Results) == 0 {
		if alt, _ := c.fetchMailsViaJWTParam(limit, offset); len(alt) > 0 {
			return alt, nil
		}
	}
	return resp.Results, nil
}

func (c *TEmailClient) fetchMailsViaJWTParam(limit, offset int) ([]temailMail, error) {
	u := fmt.Sprintf("%s/api/mails?limit=%d&offset=%d&jwt=%s", c.cfg.BaseURL, limit, offset, url.QueryEscape(c.cfg.JWT))
	body, status, err := c.request("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("获取邮件失败 (%d)", status)
	}
	var resp struct {
		Results []temailMail `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析邮件列表失败: %w", err)
	}
	return resp.Results, nil
}

func (c *TEmailClient) GetLatestMailID() (int64, error) {
	if err := c.ensureJWT(); err != nil && c.cfg.AdminPassword == "" {
		return 0, err
	}
	mails, err := c.fetchMails(1, 0)
	if err != nil {
		return 0, err
	}
	if len(mails) == 0 {
		return 0, nil
	}
	return mails[0].ID, nil
}

func (c *TEmailClient) Test() error {
	if err := c.ensureJWT(); err != nil && c.cfg.AdminPassword == "" {
		return err
	}
	_, err := c.fetchMails(1, 0)
	return err
}

// WaitForCode 轮询直到拿到验证码;只处理 ID > startMailID 的邮件
func (c *TEmailClient) WaitForCode(startMailID int64, timeoutSec, intervalSec int) (string, error) {
	if err := c.ensureJWT(); err != nil && c.cfg.AdminPassword == "" {
		return "", err
	}
	maxAttempts := timeoutSec / intervalSec
	if maxAttempts <= 0 {
		maxAttempts = 1
	}
	processed := make(map[int64]bool)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		mails, err := c.fetchMails(20, 0)
		if err != nil {
			if attempt%5 == 0 {
				log.Printf("[TEmail] 获取邮件失败: %v", err)
			}
			time.Sleep(time.Duration(intervalSec) * time.Second)
			continue
		}
		for _, mail := range mails {
			if processed[mail.ID] {
				continue
			}
			if mail.ID <= startMailID {
				processed[mail.ID] = true
				continue
			}
			if code := ExtractVerificationCode(mail.Raw); code != "" {
				log.Printf("[TEmail] 找到验证码: %s (邮件 ID %d)", code, mail.ID)
				return code, nil
			}
			processed[mail.ID] = true
		}
		if attempt < maxAttempts {
			time.Sleep(time.Duration(intervalSec) * time.Second)
		}
	}
	return "", fmt.Errorf("等待验证码超时 (%ds)", timeoutSec)
}

var (
	reAWSCode      = regexp.MustCompile(`(?i)Verification code::\s*(\d{6})`)
	reAWSCodeCN    = regexp.MustCompile(`验证码[：:][：:]\s*(\d{6})`)
	reCommonCode   = regexp.MustCompile(`(?i)Verification code:\s*(\d{6})`)
	reCodeIs       = regexp.MustCompile(`(?i)code\s+is\s+(\d{6})`)
	reChineseCode  = regexp.MustCompile(`验证码[：:]\s*(\d{6})`)
	reGenericCode  = regexp.MustCompile(`(?:^|[^\d])(\d{6})(?:[^\d]|$)`)
)

func ExtractVerificationCode(raw string) string {
	if raw == "" {
		return ""
	}
	for _, re := range []*regexp.Regexp{reAWSCode, reAWSCodeCN, reCommonCode, reCodeIs, reChineseCode} {
		if m := re.FindStringSubmatch(raw); len(m) > 1 {
			return m[1]
		}
	}
	if m := reGenericCode.FindStringSubmatch(raw); len(m) > 1 {
		return m[1]
	}
	return ""
}
