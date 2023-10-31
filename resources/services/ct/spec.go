package ct

import (
	"fmt"
	"strings"

	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/thoas/go-funk"
)

// Spec is the configuration for a list source
type Spec struct {
	// REST `$select` OData modificator, fields entity properties array
	// Wildcard selectors `*` are intentionally not supported
	// If not provided, only default fields will be fetched (ID, Created, AuthorId, Modified, EditorId)
	Select []string `json:"select"`
	// REST `$expand` OData modificator, fields entity properties array
	// When expanding an entity use selection of a nested entity property(s)
	// Optional, and in most of the cases we recommend to avoid it and
	// prefer to map nested entities to the separate tables
	Expand []string `json:"expand"`
	// Optional, an alias for the table name
	// Don't map different lists to the same table - such scenario is not supported
	Alias string `json:"alias"`

	// Custom fields mapping settings
	fieldsMapping map[string]string
}

// SetDefault sets default values for list spec
func (s *Spec) SetDefault() {
	if s.Select == nil {
		s.Select = []string{}
	}

	exclude := []string{}
	prepProps := []string{"UniqueId", "ParentList/ParentWebUrl", "ParentList/Id", "ID"}
	apndProps := []string{"Created", "AuthorId", "Modified", "EditorId"}
	defaultExpand := []string{"ParentList"}

	// Extract arrow syntax fields mapping
	s.fieldsMapping = util.GetFieldsMapping(s.Select)
	for i, field := range s.Select {
		f, _ := util.GetFieldMapping(field)
		s.Select[i] = f
	}

	s.fieldsMapping["ParentList/ParentWebUrl"] = "WebUrl"
	s.fieldsMapping["ParentList/Id"] = "ListId"

	s.Select = funk.FilterString(s.Select, func(field string) bool {
		// Disable wildcard or nested wildcard selectors
		if strings.Contains(field, "*") {
			return false
		}
		return !funk.ContainsString(util.ConcatSlice(exclude, util.ConcatSlice(prepProps, apndProps)), field)
	})

	s.Expand = util.ConcatSlice(s.Expand, defaultExpand)
	s.Select = util.ConcatSlice(prepProps, util.ConcatSlice(s.Select, apndProps))
}

// Validate validates list spec
func (s *Spec) Validate() error {
	aliases := make([]string, len(s.Select))
	for i, field := range s.Select {
		aliases[i] = util.NormalizeEntityName(field)
		if alias, ok := s.fieldsMapping[field]; ok {
			aliases[i] = util.NormalizeEntityName(alias)
		}
	}

	// All aliases should be unique, output which is not unique
	for i, alias := range aliases {
		if funk.ContainsString(aliases[i+1:], alias) {
			return fmt.Errorf("alias \"%s\" is not unique", alias)
		}
	}

	return nil
}

// GetAlias returns an alias for a content type rollup
func (*Spec) GetAlias(ctName string) string {
	return strings.ToLower("rollup_" + util.NormalizeEntityName(ctName))
}
