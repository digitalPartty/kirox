package email

import (
	"fmt"
	"log"
)

// duckTEmailAdapter 组合 DuckDuckGo 别名 + TEmail 取码,实现 TempEmailService
type duckTEmailAdapter struct {
	duckToken   string
	temailCfg   TEmailConfig
	address     string
	startMailID int64
}

func NewDuckTEmailService(duckToken string, temailCfg TEmailConfig) TempEmailService {
	return &duckTEmailAdapter{
		duckToken: duckToken,
		temailCfg: temailCfg,
	}
}

func (a *duckTEmailAdapter) Create() string {
	tc := NewTEmailClient(a.temailCfg)
	latestID, err := tc.GetLatestMailID()
	if err != nil {
		log.Printf("[DuckTEmail] 获取 TEmail 最新邮件 ID 失败: %v, 使用 0", err)
	}
	a.startMailID = latestID

	duck := NewDuckDuckGoClient(a.duckToken)
	addr, err := duck.CreateAlias()
	if err != nil {
		log.Printf("[DuckTEmail] 创建 DuckDuckGo 别名失败: %v", err)
		return ""
	}
	a.address = addr
	log.Printf("[DuckTEmail] 别名: %s (基准邮件 ID: %d)", addr, a.startMailID)
	return addr
}

func (a *duckTEmailAdapter) WaitForCode(timeoutSec, intervalSec int) (string, error) {
	if a.address == "" {
		return "", fmt.Errorf("未创建邮箱别名")
	}
	tc := NewTEmailClient(a.temailCfg)
	return tc.WaitForCode(a.startMailID, timeoutSec, intervalSec)
}

func (a *duckTEmailAdapter) GetAddress() string {
	return a.address
}

// directMailAdapter 将 DirectMailClient 包装为 TempEmailService
type directMailAdapter struct {
	cfg    DirectMailConfig
	client *DirectMailClient
}

func NewDirectMailService(cfg DirectMailConfig) TempEmailService {
	return &directMailAdapter{cfg: cfg}
}

func (a *directMailAdapter) Create() string {
	a.client = NewDirectMailClient(a.cfg)
	log.Printf("[DirectMail] 使用邮箱: %s", a.cfg.Email)
	return a.cfg.Email
}

func (a *directMailAdapter) WaitForCode(timeoutSec, intervalSec int) (string, error) {
	if a.client == nil {
		return "", fmt.Errorf("DirectMail 客户端未初始化")
	}
	return a.client.WaitForCode(timeoutSec, intervalSec)
}

func (a *directMailAdapter) GetAddress() string {
	return a.cfg.Email
}
