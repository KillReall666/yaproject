package update

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type metricsUpdate interface {
	SaveMetrics(request *model.Metrics) error
	MetricsPrint()
	GetCountMetrics(request *model.Metrics) (int64, error)
	GetFloatMetrics(response *model.Metrics) (float64, error)
}

type Handler struct {
	metricsUpdate metricsUpdate
}

func NewUpdateHandler(s metricsUpdate) *Handler {
	return &Handler{
		metricsUpdate: s,
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
			_ = h.metricsUpdate.SaveMetrics(dto)

		} else if metricsType == "gauge" {
			floatValue = floatValueConv(metricsValue)
			dto := &model.Metrics{
				Name:  metricsName,
				Gauge: &floatValue,
			}
			_ = h.metricsUpdate.SaveMetrics(dto)
		}
	}
	//	h.metricsUpdate.MetricsPrint()
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
		_ = h.metricsUpdate.SaveMetrics(dto)
	} else if metrics.MType == "gauge" {
		dto := &model.Metrics{
			Name:  metrics.ID,
			Gauge: metrics.Value,
		}
		_ = h.metricsUpdate.SaveMetrics(dto)
	}

	dto := &model.Metrics{
		Name: metrics.ID,
	}

	var metricsForGaugeRequest, metricsForCounterRequest model.MetricsJSON

	if metrics.MType == "gauge" {
		value, err1 := h.metricsUpdate.GetFloatMetrics(dto)
		metricsForGaugeRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "gauge",
			Value: handlers.Float64Ptr(value),
		}
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusNotFound)
		}
	}

	jsonGaugeData, err := json.Marshal(metricsForGaugeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Accept-Encoding", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonGaugeData)
	if err != nil {
		log.Println(err)
	}

	if metrics.MType == "counter" {
		value, err2 := h.metricsUpdate.GetCountMetrics(dto)
		metricsForCounterRequest = model.MetricsJSON{
			ID:    metrics.ID,
			MType: "counter",
			Delta: handlers.Int64Ptr(value),
		}
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
		}

		jsonCounterData, err := json.Marshal(metricsForCounterRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Accept-Encoding", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonCounterData)
		if err != nil {
			log.Println(err)
		}
		//h.metricsUpdate.MetricsPrint()
	}
}
