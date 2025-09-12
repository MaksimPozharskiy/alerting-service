package main

import (
	"alerting-service/internal/agent"
	"alerting-service/internal/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	done := make(chan struct{})

	go func() {
		defer close(done)
		agent.RuntimeAgent(ctx, client)
	}()

	select {
	case sig := <-sigint:
		logger.Log.Info("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
	case <-done:
		logger.Log.Info("Agent completed normally")
		return
	}

	<-done
	logger.Log.Info("Agent shutdown completed gracefully")
}

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
