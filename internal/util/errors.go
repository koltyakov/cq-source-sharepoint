package util

import (
	"errors"
	"strings"
)

// IsNotFound unwraps API response errors checking for 404 Not Found
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	for {
		if err.Error() == "404 Not Found" || strings.Contains(err.Error(), "System.IO.FileNotFoundException") {
			return true
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	return false
}
