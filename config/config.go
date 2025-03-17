package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// BrowserConfig 浏览器配置
type BrowserConfig struct {
	Type      string `json:"type"`      // 浏览器类型：chromium, firefox, webkit
	Headless  bool   `json:"headless"`  // 是否无头模式
	SlowMo    int    `json:"slowMo"`    // 慢动作模式，毫秒
	Maximized bool   `json:"maximized"` // 是否最大化
}

// LoginConfig 登录配置
type LoginConfig struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	URL      string `json:"url"`      // 登录URL
}

// Config 应用配置
type Config struct {
	Browser BrowserConfig `json:"browser"` // 浏览器配置
	Login   LoginConfig   `json:"login"`   // 登录配置
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Browser: BrowserConfig{
		Type:      "chromium",
		Headless:  false,
		SlowMo:    0,
		Maximized: true,
	},
	Login: LoginConfig{
		Username: "tomsmith",
		Password: "SuperSecretPassword!",
		URL:      "http://the-internet.herokuapp.com/login",
	},
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果配置文件不存在，创建默认配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 确保目录存在
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("无法创建配置目录: %w", err)
		}

		// 写入默认配置
		file, err := os.Create(configPath)
		if err != nil {
			return nil, fmt.Errorf("无法创建配置文件: %w", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(DefaultConfig); err != nil {
			return nil, fmt.Errorf("无法写入默认配置: %w", err)
		}

		return &DefaultConfig, nil
	}

	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %w", err)
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %w", err)
	}

	return config, nil
}