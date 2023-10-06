package storage

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
)

type Metrics struct {
	Count int64   `json:"count"`
	Gauge float64 `json:"gauge"`
}
type MemStorage struct {
	storage map[string]*Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]*Metrics),
	}
}

func (ms *MemStorage) CountSetter(name string, count int64) {
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Count += count
		return
	}
	ms.storage[name].Count += count
}

func (ms *MemStorage) GaugeSetter(name string, gauge float64) {
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Gauge = gauge
		return
	}
	ms.storage[name].Gauge = gauge
}

func (ms *MemStorage) GaugeGetter(key string) (float64, error) {
	_, ok := ms.storage[key]
	if !ok {
		return 0, fmt.Errorf("value with key '%s' not found", key)
	}
	return ms.storage[key].Gauge, nil
}

func (ms *MemStorage) CountGetter(key string) (int64, error) {
	_, ok := ms.storage[key]
	if !ok {
		return 0, fmt.Errorf("value with key '%s' not found", key)
	}
	return ms.storage[key].Count, nil
}

func (ms *MemStorage) GetAllMetrics() string {
	htmlPage := "Metric List\n"
	for key, metric := range ms.storage {
		if key != "PollCount" {
			htmlPage += fmt.Sprintf("%v: %v\n", html.EscapeString(key), metric.Gauge)
		} else {
			htmlPage += fmt.Sprintf("%v: %v\n", html.EscapeString(key), metric.Count)
		}
	}
	htmlPage += ""
	return htmlPage
}

func (ms *MemStorage) Print() {
	var metrics string
	for key, value := range ms.storage {
		if key != "PollCount" {
			metrics += fmt.Sprintf("%s:%v. ", key, value.Gauge)
		} else {
			metrics += fmt.Sprintf("%s:%v. ", key, value.Count)
		}
	}
	log.Println("New received metrics: ", metrics)
}

func (m *MemStorage) ToJSON() ([]byte, error) {
	return json.Marshal(m.storage)
}

func (m *MemStorage) UnmarshalJSONData(data []byte) error {
	storageData := make(map[string]json.RawMessage)

	err := json.Unmarshal(data, &storageData)
	if err != nil {
		return err
	}

	for key, value := range storageData {
		metricsData := Metrics{}
		err := json.Unmarshal(value, &metricsData)
		if err != nil {
			return err
		}

		m.storage[key] = &metricsData
	}
	return nil
}
