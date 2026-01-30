package services

import (
	"errors"
	"time"
)

func isTenDigits(s string) bool {
	if len(s) != 10 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func parseDate(raw string, endOfDay bool) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		t = t.UTC()
		if endOfDay {
			return t.Add(24*time.Hour - time.Nanosecond), nil
		}
		return t, nil
	}
	return time.Time{}, errors.New("invalid date")
}
