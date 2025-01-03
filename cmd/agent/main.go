package main

import (
	"alerting-service/internal/agent"
	"net/http"
)

func main() {
	parseFlags()

	client := &http.Client{
		Transport: &http.Transport{},
	}

	agent.RuntimeAgent(client, flagPollInterval, flagReportInterval, flagRunAddr)
}
