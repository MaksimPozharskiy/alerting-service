package main

import (
	handlers "alerting-service/internal/handlers"
	repositories "alerting-service/internal/repository"
	"alerting-service/internal/usecases"
	"fmt"
	"net/http"
)

func main() {
	storageRepository := repositories.NewStorageRepository()
	metricUsecase := usecases.NewMetricUsecase(storageRepository)
	metricsHandler := handlers.NewMetricHandler(metricUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, metricsHandler.UpdateMetric)

	fmt.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
