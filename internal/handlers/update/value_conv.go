package update

import "strconv"

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
	intVal, _ := strconv.Atoi(value)
	return int64(intVal)
}
