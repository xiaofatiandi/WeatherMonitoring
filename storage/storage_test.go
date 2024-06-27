/*
This file contains all the unit tests for storage.go
*/
package storage

import (
	"testing"
	"time"
	"weather-monitoring/logger" // Import from the module weather-monitoring

	"github.com/stretchr/testify/assert"
)

var testStorage *InMemoryStorage

func init() {
	testLogger := logger.NewLogger()
	testStorage = NewInMemoryStorage(*testLogger)
}
func TestEnrollDevice(t *testing.T) {

	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)

	assert.True(t, testStorage.IsDeviceEnrolled(deviceID), "Device should be enrolled")
}

func TestDisableDevice(t *testing.T) {

	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)

	err := testStorage.DisableDevice(deviceID)
	assert.NoError(t, err, "Disabling device should not return an error")

	assert.False(t, testStorage.IsDeviceEnrolled(deviceID), "Device should be disabled")
}

func TestEnableDevice(t *testing.T) {

	deviceID := "test-device"
	testStorage.EnrollDevice(deviceID)

	err := testStorage.DisableDevice(deviceID)
	assert.NoError(t, err, "Disabling device should not return an error")

	err = testStorage.EnableDevice(deviceID)
	assert.NoError(t, err, "Enabling device should not return an error")

	assert.True(t, testStorage.IsDeviceEnrolled(deviceID), "Device should be enabled")
}

func TestRecordTemperature(t *testing.T) {

	deviceID := "test-device"
	temperature := 25.5

	testStorage.EnrollDevice(deviceID)
	testStorage.RecordTemperature(deviceID, temperature)

	// Check if temperature is recorded for the device
	assert.NotNil(t, testStorage.temperature[deviceID], "Temperature should be recorded for the device")

}

func TestGetDailyAggregatedData(t *testing.T) {
	//clear storage
	for k := range testStorage.devices {
		delete(testStorage.devices, k)
	}

	testStorage.EnrollDevice("device1")
	testStorage.EnrollDevice("device2")
	testStorage.EnrollDevice("device3")

	testStorage.RecordTemperature("device1", 20)
	testStorage.RecordTemperature("device1", 23)
	testStorage.RecordTemperature("device1", 25)

	testStorage.RecordTemperature("device2", 10)
	testStorage.RecordTemperature("device2", 11)

	testStorage.RecordTemperature("device3", 18)

	// Get aggregated data for today
	date := time.Now()
	aggregatedData := testStorage.GetDailyAggregatedData(date)

	// Assuming more detailed checks on aggregated data
	// For brevity, assert on basic structure
	assert.NotNil(t, aggregatedData["device1"], "Aggregated data should exist for the device")
	assert.NotNil(t, aggregatedData["device2"], "Aggregated data should exist for the device")
	assert.NotNil(t, aggregatedData["device3"], "Aggregated data should exist for the device")

	assert.Equal(t, 20.0, aggregatedData["device1"].Low)
	assert.Equal(t, 25.0, aggregatedData["device1"].High)
	assert.Equal(t, 22.666666666666668, aggregatedData["device1"].Average)

	assert.Equal(t, 10.0, aggregatedData["device2"].Low)
	assert.Equal(t, 11.0, aggregatedData["device2"].High)
	assert.Equal(t, 10.5, aggregatedData["device2"].Average)

	assert.Equal(t, 18.0, aggregatedData["device3"].Low)
	assert.Equal(t, 18.0, aggregatedData["device3"].High)
	assert.Equal(t, 18.0, aggregatedData["device3"].Average)
}
