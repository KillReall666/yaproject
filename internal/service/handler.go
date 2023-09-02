package service

import (
	"github.com/KillReall666/yaproject/internal/model"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Storager interface {
	PostHandler()
}

// Get's URL as string from request
func getURL(r *http.Request) string {
	url := r.URL.String()
	return url
}

// newMetricStorage Init map
func newMetricStorage() *model.MetricsStorage {
	return &model.MetricsStorage{}
}

// valueConv Convert string from URL to int64 or float64 value, or return 0 when it's unpossible.
func valueConv(value string) any {
	for _, v := range value {
		if v == 46 {
			valForReturn, _ := strconv.ParseFloat(value, 64)
			return valForReturn
		} else {
			valForReturn, _ := strconv.Atoi(value)
			return valForReturn
		}
	}
	return 0
}

// setMetrics Put data in MetricsStorage
func setMetrics(typeOfCounter string, nameOfMetric string, value any, storage model.MetricsStorage) *model.MetricsStorage {
	storage[nameOfMetric] = &model.Metrics{}
	if typeOfCounter == "counter" {
		intVal, ok := value.(int64)
		if ok {
			storage[nameOfMetric].MetricsName = nameOfMetric
			storage[nameOfMetric].MetricsType = typeOfCounter
			storage[nameOfMetric].MetValCounter += intVal
		}
	} else {
		floatVal, ok := value.(float64)
		if ok {
			storage[nameOfMetric].MetricsName = nameOfMetric
			storage[nameOfMetric].MetricsType = typeOfCounter
			storage[nameOfMetric].MetValGauge = floatVal
		}
	}

	return &storage
}

// PostHandler Main Handler for Post-response
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusNotFound)
		return
	}

	url := getURL(r)

	urlWithoutPref, err := strings.CutPrefix(url, "/")
	if err != true {
		panic(err)
	}

	metricsString := strings.Split(urlWithoutPref, "/")
	log.Println(metricsString)

	if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
		return
	}
	numForSetMetrics := valueConv(metricsString[3])

	if metricsString[1] != "counter" && metricsString[1] != "gauge" || numForSetMetrics == 0 {
		http.Error(w, "error 400", http.StatusBadRequest)
	} else if len(metricsString) < 4 {
		http.Error(w, "error 404", http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if len(metricsString) == 4 {
		res := setMetrics(metricsString[1], metricsString[2], numForSetMetrics, *storage)
		log.Println(*res)
	}

}
