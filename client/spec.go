package client

import (
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

// Spec is the configuration for a SharePoint source
type Spec struct {
	// Gosip auth config connection params https://go.spflow.com/auth/overview
	Auth struct {
		// Auth strategy: azurecert, azurecreds, device, saml, addin, adfs, ntlm, tmg, fba
		Strategy string `json:"strategy"`
		// `creds` options are unique for different auth strategies. See more details in [Auth strategies](https://go.spflow.com/auth/strategies)
		Creds map[string]string `json:"creds"`
	} `json:"auth"`

	// A map of URIs to the list configuration
	// If no lists are provided, nothing will be fetched
	Lists map[string]ListSpec `json:"lists"`
}

// ListSpec is the configuration for a list source
type ListSpec struct {
	// REST's `$select` OData modificator, fields entity properties array
	// Wildcard selectors `*` are intentionally not supported
	// If not provided, only default fields will be fetched (ID, Created, AuthorId, Modified, EditorId)
	Select []string `json:"select"`
	// REST's `$expand` OData modificator, fields entity properties array
	// When expanding an entity use selection of a nested entity property(s)
	// Optional, and in most of the cases we recommend to avoid it and
	// prefer to map nested entities to the separate tables
	Expand []string `json:"expand"`
	// REST's `$filter` OData modificator, a filter string
	// Don't use filters for large entities which potentially can return more than 5000 in a view
	// such filtering will throttle no matter top limit is set
	Filter string `json:"filter"`
	// REST's `$top` OData modificator, a number of items to fetch per page
	// If not provided, 5000 will be used
	// In most of the cases you don't need to change this value
	// It also can't be larger than 5000 anyways
	Top int `json:"top"`
	// Optional, an alias for the table name
	// Don't map different lists to the same table - such scenariou is not supported
	Alias string `json:"alias"`
}

// SetDefaults sets default values for top level spec
func (s *Spec) SetDefaults() {
	if s.Lists == nil {
		s.Lists = make(map[string]ListSpec)
	}

	// Set default values for list specs
	for ListURI, listSpec := range s.Lists {
		listSpec.SetDefault()
		s.Lists[ListURI] = listSpec
	}
}

// SetDefault sets default values for list spec
func (l *ListSpec) SetDefault() {
	if l.Select == nil {
		l.Select = []string{}
	}

	exclude := []string{"*"}
	prepProps := []string{"ID"}
	apndProps := []string{"Created", "AuthorId", "Modified", "EditorId"}

	l.Select = funk.FilterString(l.Select, func(field string) bool {
		return !funk.ContainsString(concatSlice(exclude, concatSlice(prepProps, apndProps)), field)
	})

	l.Select = concatSlice(prepProps, concatSlice(l.Select, apndProps))
}

// Validate validates SharePoint source spec validity
func (s Spec) Validate() error {
	if s.Auth.Strategy == "" {
		return fmt.Errorf("auth.strategy is required")
	}
	if len(s.Auth.Creds) == 0 {
		return fmt.Errorf("auth.creds is required")
	}

	if len(s.Lists) == 0 {
		return fmt.Errorf("no lists configuration is provided")
	}

	// All lists should have unique aliases
	aliases := make(map[string]bool)
	for listURI, listSpec := range s.Lists {
		alias := strings.ToLower(listSpec.Alias)
		if alias == "" {
			alias = strings.ToLower(listURI)
		}
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for list \"%s\" configuration", alias, listURI)
		}
		aliases[alias] = true
	}

	return nil
}
