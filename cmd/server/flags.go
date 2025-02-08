package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	flagRunAddr         string
	flagLogLevel        string
	flagStoreInterval   int
	flagFileStoragePath string
	flagRestore         bool
)

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.IntVar(&flagStoreInterval, "i", 300, "interbal for storing data on a disk")
	flag.StringVar(&flagFileStoragePath, "f", "./backup", "path to storing file")
	flag.BoolVar(&flagRestore, "r", true, "restote or not data from file after running server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return err
		}
		flagStoreInterval = value
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		boolValue, err := strconv.ParseBool(envRestore)
		if err != nil {
			return err
		}
		flagRestore = boolValue
	}

	return nil
}
