package update

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type BatchSaveMetrics interface {
	SaveMetrics(ctx context.Context, request *model.Metrics) error
	GetCountMetrics(ctx context.Context, request *model.Metrics) (int64, error)
	GetFloatMetrics(ctx context.Context, response *model.Metrics) (float64, error)
}

type BatchHandler struct {
	BatchSaveMetrics BatchSaveMetrics
	logger           Logger
	cfg              config.RunConfig
}

func NewBatchHandler(sm BatchSaveMetrics, l Logger, c config.RunConfig) *BatchHandler {
	return &BatchHandler{
		BatchSaveMetrics: sm,
		logger:           l,
		cfg:              c,
	}
}

func (h *BatchHandler) BatchUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "поддерживает только POST запросы!", http.StatusNotFound)
		return
	}

	var metricsPack, metricsForRequestPack []model.MetricsJSON
	var metricsForRequest model.MetricsJSON

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := json.NewDecoder(r.Body).Decode(&metricsPack); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.LogInfo("ошибка при декодировании батча метрик: ", err)
		return
	}

	for _, metrics := range metricsPack {
		if metrics.MType == "counter" {
			dto := &model.Metrics{
				Name:    metrics.ID,
				Counter: metrics.Delta,
			}
			err := h.BatchSaveMetrics.SaveMetrics(ctx, dto)
			if err != nil {
				h.logger.LogInfo(err)
			}
		} else if metrics.MType == "gauge" {
			dto := &model.Metrics{
				Name:  metrics.ID,
				Gauge: metrics.Value,
			}
			err := h.BatchSaveMetrics.SaveMetrics(ctx, dto)
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
			floatVal, err := h.BatchSaveMetrics.GetFloatMetrics(ctx, dto)
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
			intVal, err := h.BatchSaveMetrics.GetCountMetrics(ctx, dto)
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
		h.logger.LogInfo("ошибка при записи jsonData: ", err)
	}
}
