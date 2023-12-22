package utils

import (
	"github.com/gosimple/slug"
	"math"
	"os"
	"path"
	"path/filepath"
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

func RenameFilesUsingSlug(dirPath string) {

	list, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, file := range list {

		if file.Name() == ".DS_Store" {
			continue
		}

		name := file.Name()

		filename := path.Base(name)
		extension := path.Ext(name)
		filenameWithoutExt := filename[:len(filename)-len(extension)]

		newName := slug.Make(filenameWithoutExt) + extension

		os.Rename(filepath.Join(dirPath, name), filepath.Join(dirPath, newName))
	}
}
