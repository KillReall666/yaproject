package handlers

func Float64Ptr(f float64) *float64 {
	return &f
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func ConvertToInt64(pointer *int64) int64 {
	return *pointer
}
func ConvertToFloat64(pointer *float64) float64 {
	return *pointer
}
