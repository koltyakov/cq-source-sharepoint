package main

import (
	"os"
	"strings"
)

type SourceSpec struct {
	Name         string     `json:"name"`
	Registry     string     `json:"registry"`
	Path         string     `json:"path"`
	Version      string     `json:"version"`
	Destinations []string   `json:"destination"`
	Spec         PluginSpec `json:"spec"`
}

type PluginSpec struct {
	Auth AuthSpec `json:"auth"`
}

type AuthSpec struct {
	Strategy string     `json:"strategy"`
	Creds    [][]string `json:"creds"`
}

func (s *SourceSpec) Marshal() []byte {
	credsArray := ""
	for _, c := range s.Spec.Auth.Creds {
		credsArray += "        " + c[0] + ": " + c[1] + "\n"
	}
	spec := strings.TrimSpace(`
kind: source
spec:
  name: ` + s.Name + `
  registry: ` + s.Registry + `
  path: ` + s.Path + `
  version: ` + s.Version + `
  destination: ["` + strings.Join(s.Destinations, `", "`) + `"]
  spec:
    auth:
      strategy: ` + s.Spec.Auth.Strategy + `
      creds:
` + credsArray + `
  `)
	return []byte(spec)
}

func (s *SourceSpec) Save(filename string) error {
	return os.WriteFile(filename, s.Marshal(), 0644)
}
