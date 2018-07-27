package util

import (
	"fmt"
	"os"
)

// FileSize returns the size of the given file
func FileSize(name string) (int64, error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return 0, fmt.Errorf("Size: %s", err)
	}
	return fileInfo.Size(), nil
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
