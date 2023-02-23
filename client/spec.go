package client

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/api"
	"github.com/thoas/go-funk"
)

type Spec struct {
	SiteURL      string `json:"site_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	// Lists to fetch, if empty all lists will be fetched
	Lists []string `json:"lists"`

	// ListFields is a map of list name to list of fields to fetch, if empty all DefaultFields will be fetched
	ListFields map[string][]string `json:"list_fields"`

	// IgnoreFields is the fields to always ignore
	IgnoreFields []string `json:"ignore_fields"`

	// pkColumn is the primary key column name, defaults to "ID"
	pkColumn string
}

func (s *Spec) SetDefaults() {
	if s.ListFields == nil {
		s.ListFields = make(map[string][]string)
	}

	if len(s.IgnoreFields) == 0 {
		s.IgnoreFields = []string{
			"ComplianceAssetId",
			"Attachments",
			"AppAuthor",
			"AppEditor",
			"ItemChildCount",
			"FolderChildCount",
		}
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
	if funk.ContainsString(s.IgnoreFields, field.InternalName) {
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
