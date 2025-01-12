package main

import (
	"alerting-service/internal/agent"
	"net/http"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	agent.RuntimeAgent(client)
}
