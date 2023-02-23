package client

import (
	"errors"
	"strings"
)

// It's not the perfect error handling but the client doesn't expose an error as object
// but this should be good enough for now
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
