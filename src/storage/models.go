// define data structure for the aggrated temparature data
package storage

type TemperatureData struct {
	Timestamp   int64
	Temperature float64
}

type AggregatedTemperatureData struct {
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Average float64 `json:"average"`
}
