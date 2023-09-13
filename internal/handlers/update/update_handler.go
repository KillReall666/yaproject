package update

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/model"
	"net/http"
	"strings"
)

type metricsSrv interface {
	SaveMetrics(request *model.Metrics) error
	GetFloatMetrics(response *model.Metrics) (float64, error)
	GetCountMetrics(request *model.Metrics) (int64, error)
	PrintForHTML() string
	MetricsPrint()
}

type Handler struct {
	service metricsSrv
}

func NewHandler(s metricsSrv) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) HTMLHandle(w http.ResponseWriter, r *http.Request) {
	htmlPage := h.service.PrintForHTML()
	fmt.Fprint(w, htmlPage)
}

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}
	url := getURL(r)
	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if !err {
		panic(err)
	}

	requestString := strings.Split(urlWithoutPref, "/")

	if len(requestString) < 3 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}

	if requestString[1] == "counter" {
		dto := &model.Metrics{
			Name: requestString[2],
		}
		value, err1 := h.service.GetCountMetrics(dto)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)
		}
	} else {
		dto := &model.Metrics{
			Name: requestString[2],
		}
		value, err2 := h.service.GetFloatMetrics(dto)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusNotFound)
		} else {
			fmt.Fprintln(w, value)
		}
	}
}

func (h *Handler) PostHandle(w http.ResponseWriter, r *http.Request) {
	var intValue int64
	var floatValue float64
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}
	url := getURL(r)

	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if !err {
		panic(err)
	}

	metricsString := strings.Split(urlWithoutPref, "/")

	if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}

	numForSetMetrics := IntValueConv(metricsString[3])

	if metricsString[1] != "counter" && metricsString[1] != "gauge" || numForSetMetrics == 0 {
		http.Error(w, "error 400", http.StatusBadRequest)
	} else if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if len(metricsString) == 4 {
		if metricsString[1] == "counter" {
			intValue = IntValueConv(metricsString[3])
			dto := &model.Metrics{
				Name:    metricsString[2],
				Counter: &intValue,
			}
			_ = h.service.SaveMetrics(dto)
		} else if metricsString[1] == "gauge" {
			floatValue = FloatValueConv(metricsString[3])
			dto := &model.Metrics{
				Name:  metricsString[2],
				Gauge: &floatValue,
			}
			_ = h.service.SaveMetrics(dto)
		}
	}
	h.service.MetricsPrint()
}
