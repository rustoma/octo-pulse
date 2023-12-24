package utils

import (
	"math"
	"os"
	"strings"
)

func IsProdDev() bool {
	if os.Getenv("APP_ENV") == "prod" {
		return true
	}
	return false
}

func Every[T comparable](slice []T, callback func(value T, index int) bool) bool {
	length := len(slice)

	for i := 0; i < length; i++ {
		value := slice[i]

		if !callback(value, i) {
			return false
		}
	}

	return true
}

func CalculateReadTime(text string) int {
	words := strings.Fields(text)
	wordCount := len(words)

	readTime := float64(wordCount) / 200.0                    // Calculate read time in minutes
	minutes := int(readTime)                                  // Extract whole minutes
	seconds := math.Round((readTime - float64(minutes)) * 60) // Calculate seconds

	// Round the time
	if seconds >= 30 {
		minutes++
	}

	return minutes
}
