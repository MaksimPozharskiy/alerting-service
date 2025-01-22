package main

import (
	handlers "alerting-service/internal/handlers"
	"alerting-service/internal/logger"
	repositories "alerting-service/internal/repository"
	"alerting-service/internal/server"
	"alerting-service/internal/usecases"

	"github.com/go-chi/chi/v5"
)

func main() {
	parseFlags()

	storageRepository := repositories.NewStorageRepository()
	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)

	server := server.NewServer(flagRunAddr)

	r := chi.NewRouter()
	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	r.Use(logger.ResponseLogger)
	r.Use(logger.RequestLogger)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", metricsHandler.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/", metricsHandler.GetMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", metricsHandler.GetAllMetrics)
	})

	err := server.Start(r)
	if err != nil {
		panic(err)
	}
}
