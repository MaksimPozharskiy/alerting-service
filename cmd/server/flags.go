package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"alerting-service/internal/config"
	"alerting-service/internal/logger"

	"go.uber.org/zap"
)

var (
	flagRunAddr            string
	flagLogLevel           string
	flagStoreInterval      int
	flagFileStoragePath    string
	flagRestore            bool
	flagDBConnectionString string
	flagHashKey            string
	flagCryptoKey          string
	flagConfigFile         string
)

func parseFlags() error {
	flag.StringVar(&flagConfigFile, "c", "", "path to config file")
	flag.StringVar(&flagConfigFile, "config", "", "path to config file")
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.IntVar(&flagStoreInterval, "i", 300, "interval for storing data on a disk")
	flag.StringVar(&flagFileStoragePath, "f", "./backup", "path to storing file")
	flag.StringVar(&flagDBConnectionString, "d", "", "connection string for postgres db")
	flag.StringVar(&flagHashKey, "k", "", "hash key string for generation signature")
	flag.StringVar(&flagCryptoKey, "crypto-key", "", "path to crypto key file")
	flag.BoolVar(&flagRestore, "r", true, "restore or not data from file after running server")
	flag.Parse()

	var serverConfig *config.ServerConfig
	configFile := flagConfigFile

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		var err error
		serverConfig, err = config.LoadServerConfig(configFile)
		if err != nil {
			logger.Log.Error("Failed to load config file", zap.Error(err))
			return err
		}
	}

	applyServerConfig(serverConfig)

	return nil
}

func applyServerConfig(serverConfig *config.ServerConfig) {
	if flagRunAddr == "" {
		if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
			flagRunAddr = envRunAddr
		} else if serverConfig != nil && serverConfig.Address != "" {
			flagRunAddr = serverConfig.Address
		}
	}

	if flagLogLevel == "" {
		if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
			flagLogLevel = envLogLevel
		} else if serverConfig != nil && serverConfig.LogLevel != "" {
			flagLogLevel = serverConfig.LogLevel
		}
	}

	if flagStoreInterval == 0 {
		if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
			if val, err := strconv.Atoi(envStoreInterval); err == nil {
				flagStoreInterval = val
			}
		} else if serverConfig != nil && serverConfig.StoreInterval != 0 {
			flagStoreInterval = int(time.Duration(serverConfig.StoreInterval).Seconds())
		}
	}

	if flagFileStoragePath == "" {
		if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
			flagFileStoragePath = envFileStoragePath
		} else if serverConfig != nil && serverConfig.StoreFile != "" {
			flagFileStoragePath = serverConfig.StoreFile
		}
	}

	if flagDBConnectionString == "" {
		if envDBConnectionString := os.Getenv("DATABASE_DSN"); envDBConnectionString != "" {
			flagDBConnectionString = envDBConnectionString
		} else if serverConfig != nil && serverConfig.DatabaseDSN != "" {
			flagDBConnectionString = serverConfig.DatabaseDSN
		}
	}

	if flagHashKey == "" {
		if envHashKey := os.Getenv("KEY"); envHashKey != "" && envHashKey != "none" {
			flagHashKey = envHashKey
		} else if serverConfig != nil && serverConfig.HashKey != "" {
			flagHashKey = serverConfig.HashKey
		}
	}

	if flagCryptoKey == "" {
		if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
			flagCryptoKey = envCryptoKey
		} else if serverConfig != nil && serverConfig.CryptoKey != "" {
			flagCryptoKey = serverConfig.CryptoKey
		}
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if boolValue, err := strconv.ParseBool(envRestore); err == nil {
			flagRestore = boolValue
		}
	} else if serverConfig != nil {
		flagRestore = serverConfig.Restore
	}
}
