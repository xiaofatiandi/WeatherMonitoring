/*
This file contains all the unit tests for handler.go
*/
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"weather-monitoring/logger"  // Import from the module weather-monitoring
	"weather-monitoring/storage" // Import from the module weather-monitoring
)

var testStorage *storage.InMemoryStorage
var testHandler *Handler

func init() {
	testLogger := logger.NewLogger()
	testStorage = storage.NewInMemoryStorage(*testLogger)
	testHandler = NewHandler(testStorage, *testLogger)
}

func TestEnrollDevice(t *testing.T) {
	req := httptest.NewRequest("POST", "/devices/test-device", nil)
	resp := httptest.NewRecorder()

	// Use Gorilla Mux to set the URL parameters
	req = mux.SetURLVars(req, map[string]string{"id": "test-device"})

	// Call the handler directly
	testHandler.EnrollDevice(resp, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestListDevices(t *testing.T) {
	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)
	testStorage.EnrollDevice(deviceID + "_2")

	// disable device requeset
	req := httptest.NewRequest("GET", "/devices", nil)
	resp := httptest.NewRecorder()

	// Call the handler directly to disable device
	testHandler.ListDevices(resp, req)

	// Check the response and device is disabled
	assert.Equal(t, http.StatusOK, resp.Code)

	var devices map[string]bool
	err := json.NewDecoder(resp.Body).Decode(&devices)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(devices))
	assert.Contains(t, devices, deviceID)
	assert.Contains(t, devices, deviceID+"_2")

}

func TestDisableDevice(t *testing.T) {
	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)

	// disable device requeset
	req := httptest.NewRequest("POST", "/devices/disable/test-device", nil)
	resp := httptest.NewRecorder()
	// Use Gorilla Mux to set the URL parameters
	req = mux.SetURLVars(req, map[string]string{"id": "test-device"})

	// Call the handler directly to disable device
	testHandler.DisableDevice(resp, req)

	// Check the response and device is disabled
	assert.Equal(t, http.StatusOK, resp.Code)
}
func TestDisableNoExistDevice(t *testing.T) {
	// Call the handler directly to disable device
	req := httptest.NewRequest("POST", "/devices/disable/no-exist-device", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "no-exist-device"})
	resp := httptest.NewRecorder()

	testHandler.DisableDevice(resp, req)

	// Check the response should return error
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestSubmitTemperature(t *testing.T) {
	deviceID := "test-device-temp"
	testStorage.EnrollDevice(deviceID)
	testStorage.EnableDevice(deviceID)

	data := map[string]interface{}{
		"device_id":   deviceID,
		"temperature": 25.5,
	}
	body, err := json.Marshal(data)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/temperature", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	// Call the handler directly
	testHandler.SubmitTemperature(resp, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.Code)

	testStorage.DisableDevice(deviceID)
}

func TestSubmitTemperatureFromDisabledDevice(t *testing.T) {
	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)
	testStorage.DisableDevice(deviceID)

	data := map[string]interface{}{
		"device_id":   deviceID,
		"temperature": 25.5,
	}
	body, err := json.Marshal(data)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/temperature", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	// Call the handler directly
	testHandler.SubmitTemperature(resp, req)

	// Check the response status code
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestGetAggregatedTemperature(t *testing.T) {
	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)
	testStorage.RecordTemperature(deviceID, 20.0)
	testStorage.RecordTemperature(deviceID, 25.0)
	testStorage.RecordTemperature(deviceID, 30.0)

	device2 := "test-device2"
	testStorage.EnrollDevice(device2)
	testStorage.RecordTemperature(device2, 10.0)
	testStorage.RecordTemperature(device2, 15.0)

	req := httptest.NewRequest("GET", "/temperature", nil)
	resp := httptest.NewRecorder()

	// Call the handler directly
	testHandler.GetAggregatedTemperature(resp, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, resp.Code)

	var aggregatedData map[string]storage.AggregatedTemperatureData
	err := json.NewDecoder(resp.Body).Decode(&aggregatedData)
	assert.NoError(t, err)
	assert.Contains(t, aggregatedData, deviceID)
	assert.Equal(t, 30.0, aggregatedData[deviceID].High)
	assert.Equal(t, 20.0, aggregatedData[deviceID].Low)
	assert.Equal(t, 25.0, aggregatedData[deviceID].Average)

	assert.Contains(t, aggregatedData, device2)
	assert.Equal(t, 15.0, aggregatedData[device2].High)
	assert.Equal(t, 10.0, aggregatedData[device2].Low)
	assert.Equal(t, 12.5, aggregatedData[device2].Average)

}
