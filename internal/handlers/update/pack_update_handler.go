package update

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type PackSaveMetrics interface {
	SaveMetrics(request *model.Metrics) error
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
}

type PackHandler struct {
	PackSaveMetrics PackSaveMetrics
	logger          Logger
	cfg             config.RunConfig
}

func NewPackHandler(sm PackSaveMetrics, l Logger, c config.RunConfig) *PackHandler {
	return &PackHandler{
		PackSaveMetrics: sm,
		logger:          l,
		cfg:             c,
	}
}

func (h *PackHandler) PackUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}

	var metricsPack, metricsForRequestPack []model.MetricsJSON
	var metricsForRequest model.MetricsJSON
	var floatVal float64
	var intVal int64
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsPack); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, metrics := range metricsPack {
		if metrics.MType == "counter" {
			dto := &model.Metrics{
				Name:    metrics.ID,
				Counter: metrics.Delta,
			}
			err = h.PackSaveMetrics.SaveMetrics(dto)
			if err != nil {
				h.logger.LogInfo(err)
			}
		} else if metrics.MType == "gauge" {
			dto := &model.Metrics{
				Name:  metrics.ID,
				Gauge: metrics.Value,
			}
			err = h.PackSaveMetrics.SaveMetrics(dto)
			if err != nil {
				h.logger.LogInfo(err)
			}
		}
	}
	for _, metrics := range metricsPack {
		dto := &model.Metrics{
			Name: metrics.ID,
		}
		if metrics.MType == "gauge" {
			floatVal, err = h.PackSaveMetrics.GetFloatMetrics(dto)
			metricsForRequest = model.MetricsJSON{
				ID:    metrics.ID,
				MType: "gauge",
				Value: handlers.Float64Ptr(floatVal),
			}
			metricsForRequestPack = append(metricsForRequestPack, metricsForRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
		} else {
			intVal, err = h.PackSaveMetrics.GetCountMetrics(dto)
			metricsForRequest = model.MetricsJSON{
				ID:    metrics.ID,
				MType: "counter",
				Delta: handlers.Int64Ptr(intVal),
			}
			metricsForRequestPack = append(metricsForRequestPack, metricsForRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
		}
	}

	jsonData, err := json.Marshal(metricsForRequestPack)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Accept-Encoding", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.LogInfo(err)
	}
}
