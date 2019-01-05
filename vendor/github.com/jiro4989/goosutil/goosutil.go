package goosutil

import (
	"os"
)

// Exists return file exists.
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
