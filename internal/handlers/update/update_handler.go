package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/yaproject/internal/config"
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=SaveMetrics

type SaveMetrics interface {
	SaveMetrics(request *model.Metrics) error
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
}

type Logger interface {
	LogInfo(args ...interface{})
}

type Handler struct {
	saveMetrics SaveMetrics
	logger      Logger
	cfg         config.RunConfig
}

func NewUpdateHandler(sm SaveMetrics, l Logger, c config.RunConfig) *Handler {
	return &Handler{
		saveMetrics: sm,
		logger:      l,
		cfg:         c,
	}
}

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}
	var intValue int64
	var floatValue float64

	metricsString := handlers.GetURL(r)

	if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}

	metricsType := metricsString[1]
	metricsName := metricsString[2]
	metricsValue := metricsString[3]

	numForSetMetrics := intValueConv(metricsValue)

	if metricsType != "counter" && metricsType != "gauge" || numForSetMetrics == 0 {
		http.Error(w, "error 400", http.StatusBadRequest)
	} else if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if len(metricsString) == 4 {
		if metricsType == "counter" {
			intValue = intValueConv(metricsValue)
			dto := &model.Metrics{
				Name:    metricsName,
				Counter: &intValue,
			}
			err := h.saveMetrics.SaveMetrics(dto)
			if err != nil {
				h.logger.LogInfo(err)
			}
		} else if metricsType == "gauge" {
			fmt.Println(metricsType, floatValue)
			floatValue = floatValueConv(metricsValue)
			dto := &model.Metrics{
				Name:  metricsName,
				Gauge: &floatValue,
			}
			err := h.saveMetrics.SaveMetrics(dto)
			if err != nil {
				h.logger.LogInfo(err)
			}
		}
	}
}

func (h *Handler) UpdateJSONMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	var metrics model.MetricsJSON

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metrics.MType == "counter" {
		dto := &model.Metrics{
			Name:    metrics.ID,
			Counter: metrics.Delta,
		}
		err = h.saveMetrics.SaveMetrics(dto)
		if err != nil {
			h.logger.LogInfo(err)
		}

	} else if metrics.MType == "gauge" {
		dto := &model.Metrics{
			Name:  metrics.ID,
			Gauge: metrics.Value,
		}
		err = h.saveMetrics.SaveMetrics(dto)
		if err != nil {
			h.logger.LogInfo(err)
		}
	}

	dto := &model.Metrics{
		Name: metrics.ID,
	}

	var metricsForRequest model.MetricsJSON
	var floatVal float64
	var intVal int64

	if metrics.MType == "gauge" {
		floatVal, err = h.saveMetrics.GetFloatMetrics(dto)
		metricsForRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "gauge",
			Value: handlers.Float64Ptr(floatVal),
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	} else {
		intVal, err = h.saveMetrics.GetCountMetrics(dto)
	}

	metricsForRequest = model.MetricsJSON{
		ID:    metrics.ID,
		MType: "counter",
		Delta: handlers.Int64Ptr(intVal),
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	jsonData, err := json.Marshal(metricsForRequest)
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
