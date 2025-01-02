package main

import (
	handlers "alerting-service/internal/handlers"
	repositories "alerting-service/internal/repository"
	"alerting-service/internal/server"
	"alerting-service/internal/usecases"

	"github.com/go-chi/chi/v5"
)

func main() {
	storageRepository := repositories.NewStorageRepository()
	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)

	server := server.NewServer("8080")

	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", metricsHandler.UpdateMetric)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", metricsHandler.GetMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", metricsHandler.GetAllMetrics)
	})

	err := server.Start(r)
	if err != nil {
		panic(err)
	}
}
