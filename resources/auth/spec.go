package auth

import (
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

// Spec is the configuration for a SharePoint auth
type Spec struct {
	// Auth strategy: azurecert, azurecreds, device, saml, addin, adfs, ntlm, tmg, fba
	Strategy string `json:"strategy"`
	// `creds` options are unique for different auth strategies. See more details in [Auth strategies](https://go.spflow.com/auth/strategies)
	Creds map[string]string `json:"creds"`
}

// Validate validates auth spec validity
func (s *Spec) Validate() error {
	docs := "https://go.spflow.com/auth/strategies"
	strategies := map[string]struct {
		fields []string
		docs   string
	}{
		"azurecert":  {fields: []string{"siteUrl", "tenantId", "clientId", "certPath", "certPass"}, docs: "https://go.spflow.com/auth/strategies/azure-certificate-auth"},
		"azurecreds": {fields: []string{"siteUrl", "tenantId", "clientId", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/azure-creds-auth"},
		"addin":      {fields: []string{"siteUrl", "clientId", "clientSecret"}, docs: "https://go.spflow.com/auth/strategies/addin"},
		"device":     {fields: []string{"siteUrl", "tenantId", "clientId"}, docs: "https://go.spflow.com/auth/strategies/azure-device-flow"},
		"saml":       {fields: []string{"siteUrl", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/saml"},
		"ntlm":       {fields: []string{"siteUrl", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/ntlm"},
		"adfs":       {fields: []string{"siteUrl", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/adfs"},
		"fba":        {fields: []string{"siteUrl", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/fba"},
		"tmg":        {fields: []string{"siteUrl", "username", "password"}, docs: "https://go.spflow.com/auth/strategies/tmg"},
	}

	strategiesList := funk.Keys(strategies).([]string)

	// Missing strategy
	if s.Strategy == "" {
		return fmt.Errorf("missing auth strategy, use one of these: %s; see more %s", strings.Join(strategiesList, ", "), docs)
	}

	// Check if strategy is supported
	if _, ok := strategies[s.Strategy]; !ok {
		return fmt.Errorf("unsupported auth strategy \"%s\", use one of these: %s; see more %s", s.Strategy, strings.Join(strategiesList, ", "), docs)
	}

	// Check if creds are provided
	if s.Creds == nil {
		return fmt.Errorf("missing auth creds for \"%s\" auth strategy, required props: %s; see more %s", s.Strategy, strings.Join(strategies[s.Strategy].fields, ", "), strategies[s.Strategy].docs)
	}

	// Check if all required fields are provided
	missedFields := []string{}
	for _, field := range strategies[s.Strategy].fields {
		if _, ok := s.Creds[field]; !ok {
			missedFields = append(missedFields, field)
		}
	}
	if len(missedFields) > 0 {
		return fmt.Errorf("missing required field(s) \"%s\" for \"%s\" auth strategy; see more %s", strings.Join(missedFields, ", "), s.Strategy, strategies[s.Strategy].docs)
	}

	return nil
}
