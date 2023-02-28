package util

import (
	"net/url"
	"strings"
)

// Retrives relative URL out of an absolute one
func GetRelativeURL(absURL string) string {
	u, _ := url.Parse(absURL)
	return u.Path
}

// Removes relative URL prefix from relative URL case insensitive
func RemoveRelativeURLPrefix(relURL, prefix string) string {
	if strings.HasPrefix(strings.ToLower(relURL), strings.ToLower(prefix)) {
		return strings.TrimPrefix(relURL, prefix)
	}
	return relURL
}
