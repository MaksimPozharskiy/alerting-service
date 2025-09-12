package main

import (
	"alerting-service/internal/compressor"
	"alerting-service/internal/crypto"
	"alerting-service/internal/db"
	handlers "alerting-service/internal/handlers"
	"alerting-service/internal/logger"
	"alerting-service/internal/metrics"
	"alerting-service/internal/observability"
	"alerting-service/internal/repository"
	"alerting-service/internal/server"
	"alerting-service/internal/signature"
	"alerting-service/internal/usecases"
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

	err := parseFlags()
	if err != nil {
		panic(err)
	}

	var storageRepository repository.StorageRepository
	var dbConn *sql.DB

	if flagDBConnectionString != "" {
		dbConn, err = db.Connect(flagDBConnectionString)
		if err != nil {
			panic(err)
		} else {
			defer dbConn.Close()
			storageRepository = repository.NewDBStorageRepository(dbConn)
		}
	} else {
		storageRepository = repository.NewMemStorageRepository()
	}

	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)
	obsHandler := observability.NewObsHandler(dbConn)

	server := server.NewServer(flagRunAddr)

	var privateKey *rsa.PrivateKey

	if flagCryptoKey != "" {
		privateKey, err = crypto.LoadPrivateKey(flagCryptoKey)
		if err != nil {
			logger.Log.Error("Failed to load private key, proceeding without decryption", zap.Error(err))
		}
	}

	r := chi.NewRouter()
	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	r.Use(crypto.DecryptionMiddleware(privateKey))
	r.Use(logger.RequestLogger)
	r.Use(logger.ResponseLogger)
	r.Use(compressor.GzipMiddleware)

	signature.SetServerHashKey(flagHashKey)
	r.Use(signature.HashMiddleware)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateMetric)
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateMetrics)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", metricsHandler.GetMetric)
	})

	r.Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateURLMetric)
	})

	r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
		r.Get("/", metricsHandler.GetURLMetric)
	})

	r.Get("/ping", obsHandler.HealthCheckDB)

	r.Route("/", func(r chi.Router) {
		r.Get("/", metricsHandler.GetAllMetrics)
	})

	backupController, err := metrics.NewBackupController(flagFileStoragePath)
	if err != nil {
		panic(err)
	}

	allMetrics, err := backupController.ReadMetrics()
	if err != nil {
		panic(err)
	}
	storageRepository.SetMetrics(allMetrics)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		logger.Log.Info("Received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down server", zap.Error(err))
		}

		allMetrics, err := storageRepository.GetMetrics()
		if err != nil {
			logger.Log.Error("Error getting metrics for backup", zap.Error(err))
		} else {
			if err := backupController.WriteMetrics(allMetrics); err != nil {
				logger.Log.Error("Error writing backup", zap.Error(err))
			} else {
				logger.Log.Info("Metrics successfully saved before shutdown")
			}
		}

		logger.Log.Info("Server shutdown completed")
		close(idleConnsClosed)
	}()

	go func() {
		if err := server.Start(r); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Error starting server", zap.Error(err))
			close(idleConnsClosed)
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(flagStoreInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := os.Truncate(flagFileStoragePath, 0); err != nil {
					logger.Log.Error("Error truncating backup file", zap.Error(err))
					continue
				}

				allMetrics, err := storageRepository.GetMetrics()
				if err != nil {
					logger.Log.Error("Error getting metrics for backup", zap.Error(err))
					continue
				}

				if err := backupController.WriteMetrics(allMetrics); err != nil {
					logger.Log.Error("Error writing backup", zap.Error(err))
				}
			case <-idleConnsClosed:
				return
			}
		}
	}()

	<-idleConnsClosed
	logger.Log.Info("Server stopped gracefully")
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
