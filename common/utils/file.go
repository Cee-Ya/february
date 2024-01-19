package utils

import "os"

func PathExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return false
}
