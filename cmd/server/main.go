package main

import (
	handlers "alerting-service/internal/handlers"
	repositories "alerting-service/internal/repository"
	"alerting-service/internal/server"
	"alerting-service/internal/usecases"
	"net/http"
)

func main() {
	storageRepository := repositories.NewStorageRepository()
	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)

	server := server.NewServer("8080")

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, metricsHandler.UpdateMetric)

	err := server.Start(mux)
	if err != nil {
		panic(err)
	}
}
