package update

import (
	"bytes"
	"encoding/json"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
	"net/http"
)

type PackSaveMetrics interface {
	SaveMetrics(request *model.Metrics) error
	SaveMetricsToDB(request *model.Metrics) error
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
	GetFloatMetricsFromDB(request *model.Metrics) (float64, error)
	GetCountMetricsFromDB(request *model.Metrics) (int64, error)
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
	var flag bool
	if h.cfg.DefaultDBConnStr == "" {
		flag = true
	}

	var buf bytes.Buffer
	var metricsPack []model.MetricsJSON
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
			if flag {
				_ = h.PackSaveMetrics.SaveMetrics(dto)
			} else {
				err = h.PackSaveMetrics.SaveMetricsToDB(dto)
				if err != nil {
					h.logger.LogInfo(err)
				}
			}
		} else if metrics.MType == "gauge" {
			dto := &model.Metrics{
				Name:  metrics.ID,
				Gauge: metrics.Value,
			}

			if flag {
				_ = h.PackSaveMetrics.SaveMetrics(dto)
			} else {
				err = h.PackSaveMetrics.SaveMetricsToDB(dto)
				if err != nil {
					h.logger.LogInfo(err)
				}
			}
		}
	}

	var metricsForRequestPack []model.MetricsJSON
	var metricsForRequest model.MetricsJSON
	var floatVal float64
	var intVal int64

	for _, metrics := range metricsPack {
		dto := &model.Metrics{
			Name: metrics.ID,
		}

		if metrics.MType == "gauge" {
			if flag {
				floatVal, err = h.PackSaveMetrics.GetFloatMetrics(dto)
			} else {
				floatVal, err = h.PackSaveMetrics.GetFloatMetricsFromDB(dto)
			}
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
			if flag {
				intVal, err = h.PackSaveMetrics.GetCountMetrics(dto)
			} else {
				intVal, err = h.PackSaveMetrics.GetCountMetricsFromDB(dto)
			}
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
