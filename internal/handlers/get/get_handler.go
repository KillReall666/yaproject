package get

import (
	"bytes"
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type metricsGet interface {
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
}

type Handler struct {
	metricsGet metricsGet
}

func NewGetHandler(s metricsGet) *Handler {
	return &Handler{
		metricsGet: s,
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
		} else {
			fmt.Fprintln(w, value)

			w.WriteHeader(http.StatusOK)

		}

	} else if metricsType == "gauge" {
		dto := &model.Metrics{
			Name: metricsName,
		}
		value, err2 := h.metricsGet.GetFloatMetrics(dto)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)

			w.WriteHeader(http.StatusOK)

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

	dto := &model.Metrics{
		Name: metrics.ID,
	}

	var metricsForRequest model.MetricsJSON

	if metrics.MType == "gauge" {
		value, err1 := h.metricsGet.GetFloatMetrics(dto)
		metricsForRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusNotFound)
		}
	} else {
		value, err2 := h.metricsGet.GetCountMetrics(dto)
		metricsForRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "counter",
			Delta: handlers.Int64Ptr(value),
		}
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
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

