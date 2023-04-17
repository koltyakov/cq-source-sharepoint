package main

import (
	"regexp"
	"strings"
)

func includes(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func getEntityID(entityName string) string {
	// Extract entity ID from entity name where entity ID is in [ENTITY_ID] brackets
	// e.g. "My List [123]" -> "123"
	re := regexp.MustCompile(`\[([^\[\]]+)\]`)
	matches := re.FindStringSubmatch(entityName)
	if len(matches) > 1 {
		return matches[1]
	}

	return entityName
}

func getEntityType(entityID string) string {
	// If starts with "0x" the entity is a content_type
	if strings.HasPrefix(entityID, "0x") {
		return "content_type"
	}
	return "list"
}
