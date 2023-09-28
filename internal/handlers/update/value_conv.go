package update

import (
	"strconv"
	"strings"
)

func floatValueConv(value string) float64 {
	var floatVal float64
	floatVal, _ = strconv.ParseFloat(value, 64)

	for _, v := range value {
		if v == 46 {
			floatVal, _ = strconv.ParseFloat(value, 64)
		}
	}

	return floatVal
}

func intValueConv(value string) int64 {
	var intVal int
	for _, v := range value {
		if v == 46 {
			parts := strings.Split(value, ".")
			intVal, _ = strconv.Atoi(parts[0])
			return int64(intVal)
		}
	}

	intVal, _ = strconv.Atoi(value)
	return int64(intVal)
}
