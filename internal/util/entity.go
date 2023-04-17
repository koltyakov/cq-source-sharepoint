package util

import (
	"strings"

	"github.com/cloudquery/plugin-sdk/v2/caser"
)

func NormalizeEntityName(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.Trim(s, "_")
	return s
}

func NormalizeEntityNameSnake(name string) string {
	csr := caser.New()
	s := csr.ToSnake(name)

	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.Trim(s, "_")
	return s
}
