package utils

import (
	"math"
	"os"
	"regexp"
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

func RemoveMultipleSpaces(text string) string {
	trimmedText := strings.TrimSpace(text)

	// Define the regular expression pattern to match consecutive whitespaces
	pattern := `\s+`

	// Compile the regular expression pattern
	reg := regexp.MustCompile(pattern)

	// Replace multiple whitespaces with a single space
	return reg.ReplaceAllString(trimmedText, " ") + "\n\n"
}

func EscapeSpecialCharacters(input string) string {
	// List of special characters that need to be escaped in a regular expression
	specialCharacters := []string{`\`, ".", "?", "+", "*", "|", "(", ")", "[", "]", "{", "}", "^", "$"}

	// Escape each special character in the input string
	for _, char := range specialCharacters {
		input = strings.ReplaceAll(input, char, "\\"+char)
	}

	return input
}
