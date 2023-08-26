package utils

import "os"

func IsProdDev() bool {
	if os.Getenv("APP_ENV") == "prod" {
		return true
	}
	return false
}
