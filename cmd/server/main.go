package main

import (
	"alerting-service/internal/compressor"
	"alerting-service/internal/db"
	handlers "alerting-service/internal/handlers"
	"alerting-service/internal/logger"
	"alerting-service/internal/metrics"
	"alerting-service/internal/observability"
	"alerting-service/internal/repository"
	"alerting-service/internal/server"
	"alerting-service/internal/usecases"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	err := parseFlags()
	if err != nil {
		panic(err)
	}

	dbConn, err := db.Connect(flagDBConnectionString)
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	storageRepository := repository.NewMemStorageRepository()
	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)
	obsHandler := observability.NewObsHandler(dbConn)

	server := server.NewServer(flagRunAddr)

	r := chi.NewRouter()
	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	r.Use(logger.ResponseLogger)
	r.Use(logger.RequestLogger)
	r.Use(compressor.GzipMiddleware)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateMetric)
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
	storageRepository.SetMetrics(allMetrics)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go gracefulShutdown(cancelFunc, server)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Duration(flagStoreInterval) * time.Second)

				if err := os.Truncate(flagFileStoragePath, 0); err != nil {
					panic(err)
				}

				allMetrics := storageRepository.GetMetrics()

				err = backupController.WriteMetrics(allMetrics)

				if err != nil {
					panic(err)
				}
			}
		}
	}()

	err = server.Start(r)
	if err != nil {
		panic(err)
	}
}

func gracefulShutdown(cancelFunc context.CancelFunc, srv server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	fmt.Println("graceful shutdown", s)

	cancelFunc() // Останавливаем фоновые горутины

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server:", err)
	}

	os.Exit(0)
}
