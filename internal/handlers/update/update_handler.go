package update

import (
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/model"
)

type metricsUpdate interface {
	SaveMetrics(request *model.Metrics) error
	MetricsPrint()
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
	//h.metricsUpdate.MetricsPrint()
}
