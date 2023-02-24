package client

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/api"
	"github.com/thoas/go-funk"
)

// ToDo: Design spec to correcpond SharePoint specifics
type Spec struct {
	// ToDo: add support for other auth strategies, it also should be in a separate `auth` section not on the top level
	// or probably introduce a connection_string like approach
	SiteURL      string `json:"site_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	// Lists to fetch, if empty all lists will be fetched
	Lists []string `json:"lists"`

	// ToDo: must be a nested props of a list, otherwise it's syntatically convinuent and not logical having spread entity configs in different places
	// ListFields is a map of list name to list of fields to fetch, if empty all DefaultFields will be fetched
	ListFields map[string][]string `json:"list_fields"`

	// Common service fields to ignore
	ignoreFields []string // no need it as a public property

	// pkColumn is the primary key column name, defaults to "ID"
	pkColumn string
}

func (s *Spec) SetDefaults() {
	if s.ListFields == nil {
		s.ListFields = make(map[string][]string)
	}

	s.ignoreFields = []string{
		"Id",
		"ComplianceAssetId",
		"Attachments",
		"AppAuthor",
		"AppEditor",
		"ItemChildCount",
		"FolderChildCount",
	}

	s.pkColumn = "ID"
}

func (s Spec) Validate() error {
	if s.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if s.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	if s.SiteURL == "" {
		return fmt.Errorf("site_url is required")
	}
	if len(s.Lists) == 0 {
		return fmt.Errorf("lists is required")
	}

	dupeLists := make(map[string]struct{}, len(s.Lists))
	for _, title := range s.Lists {
		name := normalizeName(title)
		if _, ok := dupeLists[name]; ok {
			return fmt.Errorf("found duplicate normalized list name in spec: %q (%q)", title, name)
		}
		dupeLists[name] = struct{}{}
	}

	if len(s.Lists) > 0 {
		for k := range s.ListFields {
			name := normalizeName(k)
			if _, ok := dupeLists[name]; !ok {
				return fmt.Errorf("found list_fields for unspecified list in spec: %q", k)
			}
		}
	}

	return nil
}

func (s Spec) ShouldSelectField(list string, field api.FieldInfo) bool {
	// ToDo: Only ignore until explicitly selected
	if funk.ContainsString(s.ignoreFields, field.InternalName) {
		return false
	}

	if fields := s.ListFields[list]; len(fields) > 0 {
		return funk.ContainsString(fields, field.InternalName)
	}

	// Ignore internal, hidden or computed fields
	if strings.HasPrefix(field.InternalName, "_") || field.Hidden || field.FieldTypeKind == 12 {
		return false
	}

	return true
}
