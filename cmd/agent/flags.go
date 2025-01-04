package main

import (
	"flag"
	"os"
	"strconv"
)

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port for sending server")
	flag.IntVar(&flagPollInterval, "p", 2, "how often to get metrics from runtime, seconds")
	flag.IntVar(&flagReportInterval, "r", 10, "how often to send metrics to server, seconds")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if envReportInterval, err := strconv.Atoi(envReportInterval); err == nil {
			flagReportInterval = envReportInterval
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if envPollInterval, err := strconv.Atoi(envPollInterval); err == nil {
			flagReportInterval = envPollInterval
		}
	}
}
