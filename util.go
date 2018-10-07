package tuna

import (
	"fmt"
	"strconv"
	"time"
)

func float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func stringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// fmtDuration the "hours:minutes:seconds" representation of a time.Duration.
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
