package model

type Metrics struct {
	Name    string
	Counter *int64
	Gauge   *float64
}
