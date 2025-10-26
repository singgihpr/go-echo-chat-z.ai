package config

import (
	"os"
)

type Config struct {
	ZhipuAPIKey string
	ZhipuAPIURL string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		ZhipuAPIKey: getEnv("ZHIPU_API_KEY", ""),
		ZhipuAPIURL: getEnv("ZHIPU_API_URL", "https://open.bigmodel.cn/api/paas/v3/model-api/chatglm_pro/invoke"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
