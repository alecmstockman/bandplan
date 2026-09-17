package handlers

import (
	"testing"
	"time"
)

func TestFormatTimeValue(t *testing.T) {
	formatTimeValue := funcMap["formatTimeValue"].(func(*time.Time, string) string)
	value := time.Date(2026, time.September, 17, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    *time.Time
		timeZone string
		want     string
	}{
		{name: "local time", value: &value, timeZone: "America/Chicago", want: "10:30"},
		{name: "invalid timezone fallback", value: &value, timeZone: "invalid", want: "15:30"},
		{name: "nil time", value: nil, timeZone: "UTC", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatTimeValue(test.value, test.timeZone); got != test.want {
				t.Errorf("formatTimeValue() = %q; want %q", got, test.want)
			}
		})
	}
}
