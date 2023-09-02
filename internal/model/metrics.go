package model

type Metrics struct {
	MetricsType   string
	MetricsName   string
	MetValGauge   float64
	MetValCounter int64
}

type MetricsStorage map[string]*Metrics
