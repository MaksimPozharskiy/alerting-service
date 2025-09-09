package config

import (
	"encoding/json"
	"os"
	"time"
)

type ServerConfig struct {
	Address       string        `json:"address"`
	Restore       bool          `json:"restore"`
	StoreInterval time.Duration `json:"store_interval"`
	StoreFile     string        `json:"store_file"`
	DatabaseDSN   string        `json:"database_dsn"`
	CryptoKey     string        `json:"crypto_key"`
	LogLevel      string        `json:"log_level"`
	HashKey       string        `json:"hash_key"`
}

func LoadServerConfig(filename string) (*ServerConfig, error) {
	if filename == "" {
		return &ServerConfig{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
