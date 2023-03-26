package client

import (
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/ct"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/search"
)

// Spec is the configuration for a SharePoint source
type Spec struct {
	// Gosip auth config connection params https://go.spflow.com/auth/overview
	Auth auth.Spec `json:"auth"`

	// A map of URIs to the list configuration
	// If no lists are provided, nothing will be fetched
	Lists map[string]lists.Spec `json:"lists"`

	// A map of TermSets GUIDs to the MMD configuration
	MMD map[string]mmd.Spec `json:"mmd"`

	// User profiles configuration
	Profiles profiles.Spec `json:"profiles"`

	// Search query results
	Search map[string]search.Spec `json:"search"`

	// Content types based rollup
	ContentTypes map[string]ct.Spec `json:"content_types"`
}

// SetDefaults sets default values for top level spec
func (s *Spec) SetDefaults() {
	if s.Lists == nil {
		s.Lists = make(map[string]lists.Spec)
	}

	// Set default values for list specs
	for ListURI, listSpec := range s.Lists {
		listSpec.SetDefault()
		s.Lists[ListURI] = listSpec
	}

	// Set default values for Search specs
	for searchName, searchSpec := range s.Search {
		searchSpec.SetDefault()
		s.Search[searchName] = searchSpec
	}

	// Set default values for MMD specs
	for terSetID, mmdSpec := range s.MMD {
		mmdSpec.SetDefault()
		s.MMD[terSetID] = mmdSpec
	}

	// Set default values for User Profiles spec
	s.Profiles.SetDefault()

	// Set default values for Content Types rollup specs
	for ctName, ctSpec := range s.ContentTypes {
		ctSpec.SetDefault()
		s.ContentTypes[ctName] = ctSpec
	}
}

// Validate validates SharePoint source spec validity
func (s *Spec) Validate() error {
	// Validation auth options
	if err := s.Auth.Validate(); err != nil {
		return err
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
		if err := listSpec.Validate(); err != nil {
			return fmt.Errorf("list \"%s\" configuration is invalid: %s", listURI, err)
		}
	}

	// All term sets should have unique aliases
	for terSetID, mmdSpec := range s.MMD {
		alias := strings.ToLower("mmd_" + mmdSpec.Alias)
		if mmdSpec.Alias == "" {
			alias = strings.ToLower("mmd_" + strings.ReplaceAll(terSetID, "-", ""))
		}
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for term set \"%s\" configuration", alias, terSetID)
		}
		aliases[alias] = true
	}

	// User profiles should have unique alias
	if s.Profiles.Enabled {
		alias := strings.ToLower("ups_profile")
		if s.Profiles.Alias != "" {
			alias = strings.ToLower("ups_" + s.Profiles.Alias)
		}
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for user profiles configuration", alias)
		}
		aliases[alias] = true
	}

	// Search spec validations
	for searchName, searchSpec := range s.Search {
		// Query text is required
		if searchSpec.QueryText == "" {
			return fmt.Errorf("queryText is required for search \"%s\" configuration", searchName)
		}

		// Unique alias name
		alias := strings.ToLower("search_" + util.NormalizeEntityName(searchName))
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for search \"%s\" configuration", alias, searchName)
		}
		aliases[alias] = true
		if err := searchSpec.Validate(); err != nil {
			return fmt.Errorf("search \"%s\" configuration is invalid: %s", searchName, err)
		}
	}

	// Content types based rollup spec validations
	for ctName, ctSpec := range s.ContentTypes {
		// Unique alias name
		alias := strings.ToLower("rollup_" + util.NormalizeEntityName(ctName))
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for content type \"%s\" configuration", alias, ctName)
		}
		aliases[alias] = true
		if err := ctSpec.Validate(); err != nil {
			return fmt.Errorf("content type rollup \"%s\" configuration is invalid: %s", ctName, err)
		}
	}

	// App only auth is not supported with search driven sources
	// ToDo: check other not user context auth strategies
	if s.Auth.Strategy == "addin" && (s.Profiles.Enabled || len(s.Search) > 0) {
		return fmt.Errorf("this auth strategy is not supported with search API, see more https://learn.microsoft.com/en-us/sharepoint/dev/solution-guidance/search-api-usage-sharepoint-add-in")
	}

	return nil
}

// getSpec unmarshals and validates the spec
func getSpec(src specs.Source) (*Spec, error) {
	var spec *Spec

	if err := src.UnmarshalSpec(&spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}

	spec.SetDefaults()

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return spec, nil
}
