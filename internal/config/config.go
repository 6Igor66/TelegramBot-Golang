package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Bot botSettings `json:"bot_settings"`
}

type botSettings struct {
	TelegramToken   string `json:"telegram_token" required:"true"`
	PollerTimeoutMS int    `json:"poller_timeout_ms" required:"true"`
	Debug           bool   `json:"debug"`
	Wallet          string `json:"wallet"`
}

func ReadConfigFromFile(filePath string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
		return nil, err
	}
	return &config, nil
}
