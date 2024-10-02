package utils

import (
	"time"
)

// GetTime looks at current time and returns day, month, year as Integers
func GetTime() (int, int, int) {
	currentTime := time.Now()

	parsedTime, _ := time.Parse("02-01-2006", currentTime.Format("02-01-2006"))

	day := parsedTime.Day()
	month := int(parsedTime.Month())
	year := parsedTime.Year()

	return day, month, year
}
