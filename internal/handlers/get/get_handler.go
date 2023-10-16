package get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/KillReall666/yaproject/internal/config"
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type metricsGet interface {
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
	GetFloatMetricsFromDB(request *model.Metrics) (float64, error)
	GetCountMetricsFromDB(request *model.Metrics) (int64, error)
}

type Handler struct {
	metricsGet metricsGet
	cfg        config.RunConfig
}

func NewGetHandler(s metricsGet, cfg config.RunConfig) *Handler {
	return &Handler{
		metricsGet: s,
		cfg:        cfg,
	}
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}

	metricsString := handlers.GetURL(r)

	if len(metricsString) < 3 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}

	metricsType := metricsString[1]
	metricsName := metricsString[2]

	if metricsType == "counter" {
		dto := &model.Metrics{
			Name: metricsName,
		}
		value, err1 := h.metricsGet.GetCountMetrics(dto)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusNotFound)
			return
		} else {
			fmt.Fprintln(w, value)
			return
		}

	} else if metricsType == "gauge" {
		dto := &model.Metrics{
			Name: metricsName,
		}
		value, err2 := h.metricsGet.GetFloatMetrics(dto)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)
			return
		}

	} else {
		http.Error(w, "error 404", http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}
	var metrics, metricsForRequest model.MetricsJSON
	var floatVal float64
	var intVal int64
	var flag bool
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := &model.Metrics{
		Name: metrics.ID,
	}

	if h.cfg.DefaultDBConnStr == "" {
		flag = true
	}

	if metrics.MType == "gauge" {
		if flag {
			floatVal, err = h.metricsGet.GetFloatMetrics(dto)
		} else {
			floatVal, err = h.metricsGet.GetFloatMetricsFromDB(dto)
		}

		metricsForRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "gauge",
			Value: handlers.Float64Ptr(floatVal),
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

	} else {
		if flag {
			intVal, err = h.metricsGet.GetCountMetrics(dto)
		} else {
			intVal, err = h.metricsGet.GetCountMetricsFromDB(dto)
		}

		metricsForRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "counter",
			Delta: handlers.Int64Ptr(intVal),
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	jsonData, err := json.Marshal(metricsForRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
