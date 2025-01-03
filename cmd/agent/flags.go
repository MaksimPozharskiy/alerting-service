package main

import (
	"flag"
)

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port for sending server")
	flag.IntVar(&flagPollInterval, "p", 2, "how often to get metrics from runtime, seconds")
	flag.IntVar(&flagReportInterval, "r", 10, "how often to send metrics to server, seconds")
	flag.Parse()
}
