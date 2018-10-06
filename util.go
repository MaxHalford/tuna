package tuna

import "strconv"

func float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func stringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
