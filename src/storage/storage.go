/*
- Define Storage interface to store and retrieve data
- Implement the Storage interface using in-memory storage
*/
package storage

import (
	"fmt"
	"sync"
	"time"
	"weather-monitoring/logger" // Import from the module weather-monitoring
)

// Storage is the interface for data storage
type Storage interface {
	EnrollDevice(deviceID string)
	ListDevices() map[string]bool
	DisableDevice(deviceID string) error
	EnableDevice(deviceID string) error
	IsDeviceEnrolled(deviceID string) bool
	RecordTemperature(deviceID string, temperature float64)
	GetDailyAggregatedData(date time.Time) map[string]AggregatedTemperatureData
}

// InMemoryStore is an in-memory implementation of the Storage interface
type InMemoryStorage struct {
	mu          sync.RWMutex                 // lock for memory access
	devices     map[string]bool              //store devices and if they are enabled
	temperature map[string][]TemperatureData //store devices and the temparature data reported by them
	logger      logger.Logger
}

// NewInMemoryStore creates a new InMemoryStore
func NewInMemoryStorage(logger logger.Logger) *InMemoryStorage {
	return &InMemoryStorage{
		devices:     make(map[string]bool),
		temperature: make(map[string][]TemperatureData),
		logger:      logger,
	}
}

func (s *InMemoryStorage) EnrollDevice(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.devices[deviceID] = true
}

func (s *InMemoryStorage) ListDevices() map[string]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.devices
}

func (s *InMemoryStorage) EnableDevice(deviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[deviceID]; !exists {
		s.logger.Info("Device ", deviceID, " not found")
		return fmt.Errorf("device %s not found", deviceID)
	}

	s.devices[deviceID] = true
	s.logger.Info("Device enabled:", deviceID)
	return nil
}

func (s *InMemoryStorage) DisableDevice(deviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[deviceID]; !exists {
		s.logger.Info("Device ", deviceID, " not found")
		return fmt.Errorf("device %s not found", deviceID)

	}
	s.devices[deviceID] = false
	s.logger.Info("Device disabled:", deviceID)
	return nil
}

func (s *InMemoryStorage) IsDeviceEnrolled(deviceID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.devices[deviceID]
}

func (s *InMemoryStorage) RecordTemperature(deviceID string, temperature float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.temperature[deviceID] = append(s.temperature[deviceID], TemperatureData{
		Timestamp:   time.Now().Unix(),
		Temperature: temperature,
	})
}

// Noted: this function now calculate high/low/avg per device.
// If the requirement is to calculate aggregated temperature across all devices,
// then need to chagne the return type to AggregatedTemperatureData and the calculation can also be simplified.
func (s *InMemoryStorage) GetDailyAggregatedData(date time.Time) map[string]AggregatedTemperatureData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aggregatedData := make(map[string]AggregatedTemperatureData)

	for deviceID, tempData := range s.temperature {
		var sum, high, low float64
		count := 0

		for _, data := range tempData {
			t := time.Unix(data.Timestamp, 0)
			if t.Year() == date.Year() && t.YearDay() == date.YearDay() {
				if count == 0 || data.Temperature > high {
					high = data.Temperature
				}
				if count == 0 || data.Temperature < low {
					low = data.Temperature
				}
				sum += data.Temperature
				count++
			}
		}

		if count > 0 {
			aggregatedData[deviceID] = AggregatedTemperatureData{
				High:    high,
				Low:     low,
				Average: sum / float64(count),
			}
		}
	}

	return aggregatedData
}
