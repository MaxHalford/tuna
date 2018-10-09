package tuna

import (
	"fmt"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	var testCases = []struct {
		err error
		msg string
	}{
		{
			err: ErrUnknownField{"cake"},
			msg: "no field named 'cake'",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			msg := tc.err.Error()
			if msg != tc.msg {
				t.Errorf("want: %s, got %s", tc.msg, msg)
			}
		})
	}
}
