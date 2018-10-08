package tuna

import (
	"fmt"
	"testing"
	"time"
)

func TestFmtDuration(t *testing.T) {
	var (
		testCases = []struct {
			d time.Duration
			s string
		}{
			{
				d: time.Duration(42) * time.Second,
				s: "00:00:42",
			},
			{
				d: time.Duration(42+60*42) * time.Second,
				s: "00:42:42",
			},
			{
				d: time.Duration(42+60*42+42*60*60) * time.Second,
				s: "42:42:42",
			},
			{
				d: time.Duration(42) * time.Nanosecond,
				s: "00:00:00",
			},
			{
				d: time.Duration(42+60*42+42*60*60*60) * time.Second,
				s: "2520:42:42",
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			var s = fmtDuration(tc.d)
			if s != tc.s {
				t.Errorf("Expected %s, got %s", tc.s, s)
			}
		})
	}
}
