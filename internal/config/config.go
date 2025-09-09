package config

import (
	"flag"
	"os"
	"strconv"
)

// Config holds configuration parameters for the agent.
type Config struct {
	RunAddr         string // Server run address
	LogLevel        string // Log level
	StoreInterval   int    // Interval in seconds for writing metrics to a file
	FileStoragePath string // Path to the file for storing metrics
	Restore         bool   // Whether to restore metrics from the file on startup
	PollInterval    int    // Interval for polling runtime metrics, in seconds
	ReportInterval  int    // Interval for reporting metrics to the server, in seconds
	HashKey         string // Secret key for signing metric payloads
	RateLimit       int    // Number of parallel workers for sending metrics
	CryptoKey       string // Path to the cryptographic key file
}

// GetConfig parses configuration from command-line flags and environment variables.
func GetConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "path to the public/private key file")
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
		if val, err := strconv.Atoi(envPollInterval); err == nil {
			cfg.PollInterval = val
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

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cfg.CryptoKey = envCryptoKey
	}

	return cfg
}
