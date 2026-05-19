package email

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"reg_go/internal/storage"
)

// ===== DuckDuckGo 配置管理 =====

func getDuckConfigPath() string {
	return filepath.Join(storage.GetDataDir(), "duckduckgo.dat")
}

func GetDuckDuckGoConfigs() []DuckDuckGoConfig {
	data, err := os.ReadFile(getDuckConfigPath())
	if err != nil {
		return []DuckDuckGoConfig{}
	}
	var configs []DuckDuckGoConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return []DuckDuckGoConfig{}
	}
	return configs
}

func SaveDuckDuckGoConfigs(configsJSON string) map[string]interface{} {
	var configs []DuckDuckGoConfig
	if err := json.Unmarshal([]byte(configsJSON), &configs); err != nil {
		return map[string]interface{}{"error": "配置格式错误: " + err.Error()}
	}
	for _, cfg := range configs {
		if cfg.Token == "" {
			return map[string]interface{}{"error": "Token 不能为空"}
		}
	}
	jsonData, _ := json.Marshal(configs)
	os.MkdirAll(filepath.Dir(getDuckConfigPath()), 0755)
	if err := os.WriteFile(getDuckConfigPath(), jsonData, 0600); err != nil {
		return map[string]interface{}{"error": "保存失败: " + err.Error()}
	}
	log.Printf("[DuckDuckGo] 已保存 %d 个配置", len(configs))
	return map[string]interface{}{"success": true}
}

func TestDuckDuckGoConnection(token string) map[string]interface{} {
	client := NewDuckDuckGoClient(token)
	addr, err := client.CreateAlias()
	if err != nil {
		return map[string]interface{}{"error": "测试失败: " + err.Error()}
	}
	return map[string]interface{}{"success": true, "alias": addr}
}

// ===== TEmail 配置管理 =====

func getTEmailConfigPath() string {
	return filepath.Join(storage.GetDataDir(), "temail.dat")
}

func GetTEmailConfigs() []TEmailConfig {
	data, err := os.ReadFile(getTEmailConfigPath())
	if err != nil {
		return []TEmailConfig{}
	}
	var configs []TEmailConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return []TEmailConfig{}
	}
	return configs
}

func SaveTEmailConfigs(configsJSON string) map[string]interface{} {
	var configs []TEmailConfig
	if err := json.Unmarshal([]byte(configsJSON), &configs); err != nil {
		return map[string]interface{}{"error": "配置格式错误: " + err.Error()}
	}
	for i, cfg := range configs {
		if cfg.Name == "" {
			return map[string]interface{}{"error": fmt.Sprintf("第 %d 个配置缺少名称", i+1)}
		}
		if cfg.BaseURL == "" {
			return map[string]interface{}{"error": fmt.Sprintf("配置 %s 缺少服务器地址", cfg.Name)}
		}
		if cfg.Email == "" {
			return map[string]interface{}{"error": fmt.Sprintf("配置 %s 缺少邮箱地址", cfg.Name)}
		}
		if cfg.JWT == "" && cfg.AdminPassword == "" {
			return map[string]interface{}{"error": fmt.Sprintf("配置 %s 需要 JWT 或 Admin 密码", cfg.Name)}
		}
	}
	jsonData, _ := json.Marshal(configs)
	os.MkdirAll(filepath.Dir(getTEmailConfigPath()), 0755)
	if err := os.WriteFile(getTEmailConfigPath(), jsonData, 0600); err != nil {
		return map[string]interface{}{"error": "保存失败: " + err.Error()}
	}
	log.Printf("[TEmail] 已保存 %d 个配置", len(configs))
	return map[string]interface{}{"success": true}
}

func TestTEmailConnection(configJSON string) map[string]interface{} {
	var cfg TEmailConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return map[string]interface{}{"error": "配置格式错误: " + err.Error()}
	}
	client := NewTEmailClient(cfg)
	if err := client.Test(); err != nil {
		return map[string]interface{}{"error": "连接失败: " + err.Error()}
	}
	return map[string]interface{}{"success": true}
}

// ===== DirectMail 配置管理 =====

func getDirectMailConfigPath() string {
	return filepath.Join(storage.GetDataDir(), "directmail.dat")
}

func GetDirectMailConfigs() []DirectMailConfig {
	data, err := os.ReadFile(getDirectMailConfigPath())
	if err != nil {
		return []DirectMailConfig{}
	}
	var configs []DirectMailConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return []DirectMailConfig{}
	}
	return configs
}

func SaveDirectMailConfigs(configsJSON string) map[string]interface{} {
	var configs []DirectMailConfig
	if err := json.Unmarshal([]byte(configsJSON), &configs); err != nil {
		return map[string]interface{}{"error": "配置格式错误: " + err.Error()}
	}
	for i, cfg := range configs {
		if cfg.Name == "" {
			return map[string]interface{}{"error": fmt.Sprintf("第 %d 个配置缺少名称", i+1)}
		}
		if cfg.BaseURL == "" || cfg.RefreshToken == "" || cfg.ClientID == "" || cfg.Email == "" {
			return map[string]interface{}{"error": fmt.Sprintf("配置 %s 缺少必填字段", cfg.Name)}
		}
	}
	jsonData, _ := json.Marshal(configs)
	os.MkdirAll(filepath.Dir(getDirectMailConfigPath()), 0755)
	if err := os.WriteFile(getDirectMailConfigPath(), jsonData, 0600); err != nil {
		return map[string]interface{}{"error": "保存失败: " + err.Error()}
	}
	log.Printf("[DirectMail] 已保存 %d 个配置", len(configs))
	return map[string]interface{}{"success": true}
}

func TestDirectMailConnection(configJSON string) map[string]interface{} {
	var cfg DirectMailConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return map[string]interface{}{"error": "配置格式错误: " + err.Error()}
	}
	client := NewDirectMailClient(cfg)
	if err := client.Test(); err != nil {
		return map[string]interface{}{"error": "连接失败: " + err.Error()}
	}
	return map[string]interface{}{"success": true}
}
