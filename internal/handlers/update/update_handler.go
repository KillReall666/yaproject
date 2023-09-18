package update

import (
	"fmt"
	"net/http"

	"github.com/KillReall666/yaproject/internal/model"
)

type metricsSrv interface {
	SaveMetrics(request *model.Metrics) error
	GetFloatMetrics(response *model.Metrics) (float64, error)
	GetCountMetrics(request *model.Metrics) (int64, error)
	PrintForHTML() string
	MetricsPrint()
}

type Handler struct {
	metricsSrv metricsSrv
}

func NewHandler(s metricsSrv) *Handler {
	return &Handler{
		metricsSrv: s,
	}
}

func (h *Handler) HTMLOutput(w http.ResponseWriter, r *http.Request) {
	htmlPage := h.metricsSrv.PrintForHTML()
	fmt.Fprint(w, htmlPage)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}

	metricsString := getURL(r)

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
		value, err1 := h.metricsSrv.GetCountMetrics(dto)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)
		}

	} else if metricsType == "gauge" {
		dto := &model.Metrics{
			Name: metricsName,
		}
		value, err2 := h.metricsSrv.GetFloatMetrics(dto)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)
		}

	} else {
		http.Error(w, "error 404", http.StatusBadRequest)
		return
	}
}

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}

	var intValue int64
	var floatValue float64

	metricsString := getURL(r)

	if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}

	metricsType := metricsString[1]
	metricsName := metricsString[2]
	metricsValue := metricsString[3]

	numForSetMetrics := IntValueConv(metricsValue)

	if metricsType != "counter" && metricsType != "gauge" || numForSetMetrics == 0 {
		http.Error(w, "error 400", http.StatusBadRequest)
	} else if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if len(metricsString) == 4 {
		if metricsType == "counter" {
			intValue = IntValueConv(metricsValue)
			dto := &model.Metrics{
				Name:    metricsName,
				Counter: &intValue,
			}
			_ = h.metricsSrv.SaveMetrics(dto)

		} else if metricsType == "gauge" {
			floatValue = FloatValueConv(metricsValue)
			dto := &model.Metrics{
				Name:  metricsName,
				Gauge: &floatValue,
			}
			_ = h.metricsSrv.SaveMetrics(dto)
		}
	}
	//h.metricsSrv.MetricsPrint()
}
