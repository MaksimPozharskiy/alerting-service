package main

import (
	"alerting-service/internal/agent"
	"context"
	"fmt"
	"net/http"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

	client := &http.Client{
		Transport: &http.Transport{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	agent.RuntimeAgent(ctx, client)
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
