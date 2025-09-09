package config

import (
	"encoding/json"
	"os"
	"time"
)

type AgentConfig struct {
	Address        string        `json:"address"`
	ReportInterval time.Duration `json:"report_interval"`
	PollInterval   time.Duration `json:"poll_interval"`
	CryptoKey      string        `json:"crypto_key"`
	HashKey        string        `json:"hash_key"`
	RateLimit      int           `json:"rate_limit"`
}

func LoadAgentConfig(filename string) (*AgentConfig, error) {
	if filename == "" {
		return &AgentConfig{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
