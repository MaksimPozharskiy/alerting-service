package config

import (
	"os"
	"testing"
)

func TestGetConfig_FromEnv(t *testing.T) {
	_ = os.Setenv("ADDRESS", "127.0.0.1:9999")
	_ = os.Setenv("REPORT_INTERVAL", "15")
	_ = os.Setenv("POLL_INTERVAL", "5")
	_ = os.Setenv("KEY", "secret")
	_ = os.Setenv("RATE_LIMIT", "10")

	cfg := GetConfig()

	if cfg.RunAddr != "127.0.0.1:9999" {
		t.Errorf("RunAddr = %s; want 127.0.0.1:9999", cfg.RunAddr)
	}
	if cfg.ReportInterval != 15 {
		t.Errorf("ReportInterval = %d; want 15", cfg.ReportInterval)
	}
	if cfg.PollInterval != 5 {
		t.Errorf("PollInterval = %d; want 5", cfg.PollInterval)
	}
	if cfg.HashKey != "secret" {
		t.Errorf("HashKey = %s; want secret", cfg.HashKey)
	}
	if cfg.RateLimit != 10 {
		t.Errorf("RateLimit = %d; want 10", cfg.RateLimit)
	}
}
