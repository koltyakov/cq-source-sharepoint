package util

import "strings"

func GetFieldsMapping(fields []string) map[string]string {
	fieldsMapping := map[string]string{}
	for _, field := range fields {
		f, m := GetFieldMapping(field)
		if len(m) > 0 {
			fieldsMapping[f] = m
		}
	}
	return fieldsMapping
}

func GetFieldMapping(field string) (string, string) {
	spl := strings.Split(field, "->")
	if len(spl) > 1 {
		return strings.Trim(spl[0], " "), strings.Trim(spl[1], " ")
	}
	return field, ""
}
