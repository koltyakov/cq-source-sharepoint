package util

import (
	"reflect"
	"strings"
)

// Extracts value from map by property path
func GetRespValByProp(val map[string]any, propPath string) any {
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
