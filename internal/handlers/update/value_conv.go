package update

import (
	"fmt"
	"strconv"
	"strings"
)

func FloatValueConv(value string) float64 {
	var floatVal float64
	for _, v := range value {
		if v == 46 {
			floatVal, _ = strconv.ParseFloat(value, 64)
		}
	}
	return floatVal
}

func IntValueConv(value string) int64 {
	var intVal int
	for _, v := range value {
		if v == 46 {
			parts := strings.Split(value, ".")
			fmt.Println(parts[0])
			intVal, _ = strconv.Atoi(parts[0])
			return int64(intVal)
		}
	}

	intVal, _ = strconv.Atoi(value)
	return int64(intVal)
}
