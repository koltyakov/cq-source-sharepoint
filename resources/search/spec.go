package search

import (
	"fmt"
	"strings"

	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/thoas/go-funk"
)

// Spec is the configuration for Search source
type Spec struct {
	QueryText        string   `json:"query_text"`        // Required
	TrimDuplicates   bool     `json:"trim_duplicates"`   // Optional, default is false
	SourceID         string   `json:"source_id"`         // Optional, default is empty
	SelectProperties []string `json:"select_properties"` // Optional, default is empty array

	// Custom fields mapping settings
	fieldsMapping map[string]string
}

// SetDefault sets default values for list spec
func (s *Spec) SetDefault() {
	if s.QueryText == "" {
		s.QueryText = "*"
	}

	// Extract arrow syntax fields mapping
	s.fieldsMapping = util.GetFieldsMapping(s.SelectProperties)
	for i, field := range s.SelectProperties {
		f, _ := util.GetFieldMapping(field)
		s.SelectProperties[i] = f
	}
}

// Validate validates search spec
func (s *Spec) Validate() error {
	defaultFields := []string{"Id", "Title"}
	selectProps := funk.FilterString(s.SelectProperties, func(field string) bool {
		return !funk.Contains(defaultFields, field)
	})
	selectProps = append(defaultFields, selectProps...)

	aliases := make([]string, len(selectProps))
	for i, field := range selectProps {
		aliases[i] = util.NormalizeEntityName(field)
		if alias, ok := s.fieldsMapping[field]; ok {
			aliases[i] = util.NormalizeEntityName(alias)
		}
	}

	// Can't use aliase for DocId
	if _, ok := s.fieldsMapping["DocId"]; ok {
		return fmt.Errorf("can't use alias for DocId, it's always \"id\"")
	}

	// All aliases should be unique, output which is not unique
	for i, alias := range aliases {
		if funk.ContainsString(aliases[i+1:], alias) {
			return fmt.Errorf("alias \"%s\" is not unique", alias)
		}
	}

	return nil
}

// GetAlias returns an alias for the list
func (s *Spec) GetAlias(searchName string) string {
	return strings.ToLower("search_" + util.NormalizeEntityName(searchName))
}
