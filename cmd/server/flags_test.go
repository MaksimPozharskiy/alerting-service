package main

import (
	"os"
	"testing"
)

func TestParseFlags_FromEnv(t *testing.T) {
	os.Setenv("ADDRESS", "127.0.0.1:9999")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("STORE_INTERVAL", "600")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/metrics")
	os.Setenv("DATABASE_DSN", "postgres://test:test@localhost:5432/testdb?sslmode=disable")
	os.Setenv("KEY", "supersecret")
	os.Setenv("RESTORE", "false")

	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("KEY")
		os.Unsetenv("RESTORE")
	}()

	err := parseFlags()
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if flagRunAddr != "127.0.0.1:9999" {
		t.Errorf("flagRunAddr = %s; want 127.0.0.1:9999", flagRunAddr)
	}
	if flagLogLevel != "debug" {
		t.Errorf("flagLogLevel = %s; want debug", flagLogLevel)
	}
	if flagStoreInterval != 600 {
		t.Errorf("flagStoreInterval = %d; want 600", flagStoreInterval)
	}
	if flagFileStoragePath != "/tmp/metrics" {
		t.Errorf("flagFileStoragePath = %s; want /tmp/metrics", flagFileStoragePath)
	}
	if flagDBConnectionString != "postgres://test:test@localhost:5432/testdb?sslmode=disable" {
		t.Errorf("flagDBConnectionString = %s; want correct DSN", flagDBConnectionString)
	}
	if flagHashKey != "supersecret" {
		t.Errorf("flagHashKey = %s; want supersecret", flagHashKey)
	}
	if flagRestore != false {
		t.Errorf("flagRestore = %v; want false", flagRestore)
	}
}
