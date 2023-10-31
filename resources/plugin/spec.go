package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/ct"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/search"
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

	// App only auth is not supported with search driven sources
	// ToDo: check other not user context auth strategies
	if s.Auth.Strategy == "addin" && (s.Profiles.Enabled || len(s.Search) > 0) {
		return fmt.Errorf("this auth strategy is not supported with search API, see more https://learn.microsoft.com/en-us/sharepoint/dev/solution-guidance/search-api-usage-sharepoint-add-in")
	}

	if err := s.validateAliases(); err != nil {
		return err
	}

	if err := s.validateLists(); err != nil {
		return err
	}

	if err := s.validateMMD(); err != nil {
		return err
	}

	if err := s.validateProfiles(); err != nil {
		return err
	}

	if err := s.validateSearch(); err != nil {
		return err
	}

	return s.validateContentTypes()
}

func (s *Spec) validateAliases() error {
	aliases := make(map[string]bool)

	for listURI, listSpec := range s.Lists {
		alias := listSpec.GetAlias(listURI)
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for list \"%s\" configuration", alias, listURI)
		}
		aliases[alias] = true
	}

	for terSetID, mmdSpec := range s.MMD {
		alias := mmdSpec.GetAlias(terSetID)
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for term set \"%s\" configuration", alias, terSetID)
		}
		aliases[alias] = true
	}

	if s.Profiles.Enabled {
		alias := s.Profiles.GetAlias()
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for user profiles configuration", alias)
		}
		aliases[alias] = true
	}

	for searchName, searchSpec := range s.Search {
		alias := searchSpec.GetAlias(searchName)
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for search \"%s\" configuration", alias, searchName)
		}
		aliases[alias] = true
	}

	for ctName, ctSpec := range s.ContentTypes {
		alias := ctSpec.GetAlias(ctName)
		if _, ok := aliases[alias]; ok {
			return fmt.Errorf("duplicate alias \"%s\" for content type \"%s\" configuration", alias, ctName)
		}
		aliases[alias] = true
	}

	return nil
}

func (s *Spec) validateLists() error {
	for listURI, listSpec := range s.Lists {
		if err := listSpec.Validate(); err != nil {
			return fmt.Errorf("list \"%s\" configuration is invalid: %s", listURI, err)
		}
	}
	return nil
}

func (s *Spec) validateMMD() error {
	for terSetID, mmdSpec := range s.MMD {
		if err := mmdSpec.Validate(); err != nil {
			return fmt.Errorf("term set \"%s\" configuration is invalid: %s", terSetID, err)
		}
	}
	return nil
}

func (s *Spec) validateProfiles() error {
	if s.Profiles.Enabled {
		if err := s.Profiles.Validate(); err != nil {
			return fmt.Errorf("user profiles configuration is invalid: %s", err)
		}
	}
	return nil
}

func (s *Spec) validateSearch() error {
	// Search spec validations
	for searchName, searchSpec := range s.Search {
		// Query text is required
		if searchSpec.QueryText == "" {
			return fmt.Errorf("queryText is required for search \"%s\" configuration", searchName)
		}

		// Validate search spec
		if err := searchSpec.Validate(); err != nil {
			return fmt.Errorf("search \"%s\" configuration is invalid: %s", searchName, err)
		}
	}
	return nil
}

func (s *Spec) validateContentTypes() error {
	for ctName, ctSpec := range s.ContentTypes {
		if err := ctSpec.Validate(); err != nil {
			return fmt.Errorf("content type rollup \"%s\" configuration is invalid: %s", ctName, err)
		}
	}
	return nil
}

// getSpec unmarshals and validates the spec
func getSpec(src []byte) (*Spec, error) {
	var spec *Spec

	if err := json.Unmarshal(src, &spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}

	spec.SetDefaults()

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return spec, nil
}
