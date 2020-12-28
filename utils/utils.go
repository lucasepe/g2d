package utils

import (
	"os"
)

// Mkdir creates a new directory
func Mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return nil
}

// Exists check if a file exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
