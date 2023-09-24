package get

import (
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
