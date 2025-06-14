package config

import (
	"flag"
	"testing"
)

func TestGetConfigEnvPollInterval(t *testing.T) {
	t.Setenv("POLL_INTERVAL", "5")
	t.Setenv("REPORT_INTERVAL", "9")

	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
	cfg := GetConfig()

	if cfg.PollInterval != 5 {
		t.Errorf("expected PollInterval=5, got %d", cfg.PollInterval)
	}
	if cfg.ReportInterval != 9 {
		t.Errorf("expected ReportInterval=9, got %d", cfg.ReportInterval)
	}
}
