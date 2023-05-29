package utils

import (
	"regexp"
	"time"
)

func CompareDates(date1, date2 time.Time) int {
	if date1.After(date2) {
		return 1
	} else if date1.Before(date2) || date1.Equal(date2) {
		return -1
	} else {
		return 0
	}
}

func ValidateDate(date string) bool {
	pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	return pattern.MatchString(date)
}
