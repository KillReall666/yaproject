package update

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/model"
	"net/http"
	"strings"
)

type Service interface {
	SaveMetrics(request *model.Metrics) error
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
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
	fmt.Println(url)
	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if err != true {
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
				Name:    metricsString[1],
				Counter: &intValue,
			}
			_ = h.service.SaveMetrics(dto)
			fmt.Println(dto.Name, *dto.Counter)
		} else if metricsString[1] == "gauge" {
			floatValue = FloatValueConv(metricsString[3])
			dto := &model.Metrics{
				Name:  metricsString[2],
				Gauge: &floatValue,
			}
			_ = h.service.SaveMetrics(dto)
			fmt.Println(dto.Name, *dto.Gauge)
		}
	}

}
