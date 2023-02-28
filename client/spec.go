package client

import (
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/mmd"
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
