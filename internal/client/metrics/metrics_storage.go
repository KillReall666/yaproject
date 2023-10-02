package metrics


import (
	"math/rand"
	"runtime"
)

type GaugeMetricsGetter struct {
	Gauge        map[string]float64
	Counter      map[string]int64
	GaugeStorage map[string]string
}

func NewGaugeMetricsStorage() *GaugeMetricsGetter {
	return &GaugeMetricsGetter{
		Gauge:        make(map[string]float64),
		GaugeStorage: make(map[string]string),
	}
}

func (gms *GaugeMetricsGetter) UpdateMetrics() map[string]string {
	v := runtime.MemStats{}
	runtime.ReadMemStats(&v)

	gms.Gauge["Alloc"] = float64(v.Alloc)
	gms.Gauge["BuckHashSys"] = float64(v.BuckHashSys)
	gms.Gauge["Frees"] = float64(v.Frees)
	gms.Gauge["GCCPUFraction"] = v.GCCPUFraction
	gms.Gauge["GCSys"] = float64(v.GCSys)
	gms.Gauge["HeapAlloc"] = float64(v.HeapAlloc)
	gms.Gauge["HeapIdle"] = float64(v.HeapIdle)
	gms.Gauge["HeapInuse"] = float64(v.HeapInuse)
	gms.Gauge["HeapObjects"] = float64(v.HeapObjects)
	gms.Gauge["HeapReleased"] = float64(v.HeapReleased)
	gms.Gauge["HeapSys"] = float64(v.HeapSys)
	gms.Gauge["LastGC"] = float64(v.LastGC)
	gms.Gauge["Lookups"] = float64(v.Lookups)
	gms.Gauge["MCacheInuse"] = float64(v.MCacheInuse)
	gms.Gauge["MCacheSys"] = float64(v.MCacheSys)
	gms.Gauge["MSpanInuse"] = float64(v.MSpanInuse)
	gms.Gauge["MSpanSys"] = float64(v.MSpanSys)
	gms.Gauge["Mallocs"] = float64(v.Mallocs)
	gms.Gauge["NextGC"] = float64(v.NextGC)
	gms.Gauge["NumForcedGC"] = float64(v.NumForcedGC)
	gms.Gauge["NumGC"] = float64(v.NumGC)
	gms.Gauge["OtherSys"] = float64(v.OtherSys)
	gms.Gauge["PauseTotalNs"] = float64(v.PauseTotalNs)
	gms.Gauge["StackInuse"] = float64(v.StackInuse)
	gms.Gauge["StackSys"] = float64(v.StackSys)
	gms.Gauge["Sys"] = float64(v.Sys)
	gms.Gauge["TotalAlloc"] = float64(v.TotalAlloc)

	gms.Gauge["RandomValue"] = rand.Float64()

	return gms.GaugeStorage
}
