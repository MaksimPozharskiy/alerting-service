package main

import (
	"alerting-service/internal/agent"
	"context"
	"net/http"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	agent.RuntimeAgent(ctx, client)
}
