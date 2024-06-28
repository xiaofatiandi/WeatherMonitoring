/*
Web request handler class
*/
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"weather-monitoring/logger"  // Import from the module weather-monitoring
	"weather-monitoring/storage" // Import from the module weather-monitoring

	"github.com/gorilla/mux"
)

type Handler struct {
	storage storage.Storage
	logger  logger.Logger
}

func NewHandler(storage storage.Storage, logger logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}

func (h *Handler) EnrollDevice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("call handler.EnrollDevice()")
	vars := mux.Vars(r)
	deviceID := vars["id"]
	h.storage.EnrollDevice(deviceID)
	json.NewEncoder(w).Encode(map[string]string{"device_id": deviceID})
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Call handler.ListDevices()")
	devices := h.storage.ListDevices()
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) EnableDevice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("call handler.EnableDevice()")
	vars := mux.Vars(r)
	deviceID := vars["id"]

	err := h.storage.EnableDevice(deviceID)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	h.logger.Info("Enabled device:", deviceID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DisableDevice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("call handler.DisableDevice()")
	vars := mux.Vars(r)
	deviceID := vars["id"]

	err := h.storage.DisableDevice(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info("Disabled device:", deviceID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SubmitTemperature(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("call handler.SubmitTemperature()")
	var req struct {
		DeviceID    string  `json:"device_id"`
		Temperature float64 `json:"temperature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.storage.IsDeviceEnrolled(req.DeviceID) {
		http.Error(w, "device not enrolled", http.StatusForbidden)
		return
	}

	h.storage.RecordTemperature(req.DeviceID, req.Temperature)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAggregatedTemperature(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("call handler.GetAggregatedTemperature()")
	aggregatedData := h.storage.GetDailyAggregatedData(time.Now())
	json.NewEncoder(w).Encode(aggregatedData)
}
