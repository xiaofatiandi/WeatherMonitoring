package main

import (
	"net/http"

	"weather-monitoring/handler"
	"weather-monitoring/logger"  // Import from the module weather-monitoring
	"weather-monitoring/storage" // Import from the module weather-monitoring

	"github.com/gorilla/mux"
)

func main() {
	// Initialize storage and logger
	logger := logger.NewLogger()

	// I used the interface for storage, so storage can be initialized using other implementations as well
	// e.g., if in the future if we want to use database storage rather than in-memory storage
	// we can just change the below line of code to replace the InMemoryStorage with DatabaseStorage class object
	storage := storage.NewInMemoryStorage(*logger)

	// Initialize handler with storage
	handler := handler.NewHandler(storage, *logger)

	// Initialize router
	router := mux.NewRouter()

	// Define routes and attach handler methods
	router.HandleFunc("/devices/{id}", handler.EnrollDevice).Methods("POST")
	router.HandleFunc("/devices", handler.ListDevices).Methods("GET")
	router.HandleFunc("/device/enable/{id}", handler.EnableDevice).Methods("POST")
	router.HandleFunc("/device/disable/{id}", handler.DisableDevice).Methods("POST")
	router.HandleFunc("/temperature", handler.SubmitTemperature).Methods("POST")
	router.HandleFunc("/temperature/aggregated", handler.GetAggregatedTemperature).Methods("GET")

	// Start server
	port := ":8000"
	logger.Info("Server started at", port)
	logger.Error(http.ListenAndServe(port, router))
}
