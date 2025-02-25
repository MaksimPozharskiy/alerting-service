package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	RunAddr         string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	PollInterval    int
	ReportInterval  int
	HashKey         string
	RateLimit       int
}

func GetConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "how often to get metrics from runtime, seconds")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "how often to send metrics to server, seconds")
	flag.IntVar(&cfg.RateLimit, "l", 3, "count of workers for sending metrics")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key string for generation signature")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.RunAddr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if envReportInterval, err := strconv.Atoi(envReportInterval); err == nil {
			cfg.ReportInterval = envReportInterval
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if envPollInterval, err := strconv.Atoi(envPollInterval); err == nil {
			cfg.ReportInterval = envPollInterval
		}
	}

	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		cfg.HashKey = envHashKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if envRateLimit, err := strconv.Atoi(envRateLimit); err == nil {
			cfg.RateLimit = envRateLimit
		}
	}

	return cfg
}
