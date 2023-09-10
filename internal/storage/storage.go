package storage

import (
	"fmt"
	"html"
	"net/http"
)

type Metrics struct {
	Count int64
	Gauge float64
}
type MemStorage struct {
	storage map[string]*Metrics
	//mtx     *sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]*Metrics),
	}
}

func (ms *MemStorage) CountSetter(name string, count int64) {
	//ms.mtx.Lock()
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Count += count
		return
	}
	ms.storage[name].Count += count
	//ms.mtx.Unlock()
}

func (ms *MemStorage) GaugeSetter(name string, gauge float64) {
	//ms.mtx.Lock()
	_, ok := ms.storage[name]
	if !ok {
		ms.storage[name] = &Metrics{}
		ms.storage[name].Gauge = gauge
		return
	}
	ms.storage[name].Gauge = gauge
	//ms.mtx.Unlock()
}

func (ms *MemStorage) GaugeGetter(key string) (float64, error) {
	_, ok := ms.storage[key]
	if !ok {
		return 0, fmt.Errorf("Значение с ключом '%s' не найдено", key)
	}
	return ms.storage[key].Gauge, nil
}

func (ms *MemStorage) CountGetter(key string) (int64, error) {
	_, ok := ms.storage[key]
	if !ok {
		return 0, fmt.Errorf("Значение с ключом '%s' не найдено", key)
	}
	return ms.storage[key].Count, nil

}

func (ms *MemStorage) GetAllMetrics(w http.ResponseWriter) {
	htmls := "Metric List\n"
	for key, metric := range ms.storage {
		htmls += fmt.Sprintf("%v: %v (%v)\n", html.EscapeString(key), metric.Count, metric.Gauge)
	}
	htmls += ""
	fmt.Fprint(w, htmls)
}
