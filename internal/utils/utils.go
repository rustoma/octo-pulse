package utils

import "os"

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
