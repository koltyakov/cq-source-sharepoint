package client

import (
	"net/url"
	"reflect"
	"strings"
)

// Retrives relative URL out of an absolute one
func getRelativeURL(absURL string) string {
	u, _ := url.Parse(absURL)
	return u.Path
}

// Removes relative URL prefix from relative URL case insensitive
func removeRelativeURLPrefix(relURL, prefix string) string {
	if strings.HasPrefix(strings.ToLower(relURL), strings.ToLower(prefix)) {
		return strings.TrimPrefix(relURL, prefix)
	}
	return relURL
}

// Extracts value from map by property path
func getRespValByProp(val map[string]any, propPath string) any {
	if val == nil {
		return nil
	}

	if strings.Contains(propPath, "/") {
		parts := strings.Split(propPath, "/")
		for _, part := range parts {
			v := val[part]
			if reflect.TypeOf(v).Kind() != reflect.Map {
				return v
			}
			val = v.(map[string]any)
		}
		return val
	}

	return val[propPath]
}
