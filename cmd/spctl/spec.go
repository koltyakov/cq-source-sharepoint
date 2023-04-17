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
	Auth         AuthSpec          `json:"auth"`
	Lists        []ListConf        `json:"lists"`
	ContentTypes []ContentTypeConf `json:"content_types"`
	MMD          []MMDConf         `json:"mmd"`
}

type AuthSpec struct {
	Strategy string     `json:"strategy"`
	Creds    [][]string `json:"creds"`
}

func (s *SourceSpec) Marshal() []byte {
	spec := `
kind: source
spec:
  name: ` + s.Name + `
  registry: ` + s.Registry + `
  path: ` + s.Path + `
  version: ` + s.Version + `
  destination: ["` + strings.Join(s.Destinations, `", "`) + `"]
  spec:
`

	if len(s.Spec.Auth.Creds) > 0 {
		spec += marshalAuth(s.Spec.Auth)
	}

	if len(s.Spec.Lists) > 0 {
		spec += marshalLists(s.Spec.Lists)
	}

	if len(s.Spec.ContentTypes) > 0 {
		spec += marshalContentTypes(s.Spec.ContentTypes)
	}

	if len(s.Spec.MMD) > 0 {
		spec += marshalMMD(s.Spec.MMD)
	}

	return []byte(strings.TrimSpace(spec))
}

func (s *SourceSpec) Save(filename string) error {
	return os.WriteFile(filename, s.Marshal(), 0644)
}

func marshalAuth(authSpec AuthSpec) string {
	res := "    auth:\n"
	res += "      strategy: " + authSpec.Strategy + "\n"
	res += "      creds:\n"
	for _, c := range authSpec.Creds {
		res += "        " + c[0] + ": " + c[1] + "\n"
	}
	return res
}

func marshalLists(listsSpec []ListConf) string {
	res := "    lists:\n"
	for _, list := range listsSpec {
		res += "      " + list.ID + ":\n"
		res += "        select:\n"
		for _, field := range list.Spec.Select {
			res += "          - " + field + "\n"
		}
		if len(list.Spec.Expand) > 0 {
			res += "        expand:\n"
			for _, field := range list.Spec.Expand {
				res += "          - " + field + "\n"
			}
		}
		if len(list.Spec.Filter) > 0 {
			res += "        filter: \"" + list.Spec.Filter + "\"\n"
		}
		if len(list.Spec.Alias) > 0 {
			res += "        alias: \"" + list.Spec.Alias + "\"\n"
		}
	}
	return res
}

func marshalContentTypes(ctSpec []ContentTypeConf) string {
	res := "    content_types:\n"
	for _, ct := range ctSpec {
		res += "      " + ct.ID + ":\n"
		res += "        select:\n"
		for _, field := range ct.Spec.Select {
			res += "          - " + field + "\n"
		}
		if len(ct.Spec.Expand) > 0 {
			res += "        expand:\n"
			for _, field := range ct.Spec.Expand {
				res += "          - " + field + "\n"
			}
		}
		if len(ct.Spec.Alias) > 0 {
			res += "        alias: \"" + ct.Spec.Alias + "\"\n"
		}
	}
	return res
}

func marshalMMD(mmdSpec []MMDConf) string {
	res := "    mmd:\n"
	for _, mmd := range mmdSpec {
		res += "      " + mmd.ID + ":\n"
		res += "        alias: \"" + mmd.Spec.Alias + "\"\n"
	}
	return res
}
