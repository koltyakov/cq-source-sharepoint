package mmd

import "strings"

// Spec is the configuration for MMD term set source
type Spec struct {
	// Optional, an alias for the table name
	// Don't map different term sets to the same table - such scenario is not supported
	Alias string `json:"alias"`
}

// SetDefault sets default values for MMD spec
func (*Spec) SetDefault() {
	// Default values
}

// Validate validates MMD spec validity
func (*Spec) Validate() error {
	// Nothing to validate
	return nil
}

// GetAlias returns the alias for the term set
func (s *Spec) GetAlias(terSetID string) string {
	if s.Alias == "" {
		return strings.ToLower("mmd_" + strings.ReplaceAll(terSetID, "-", ""))
	}
	return strings.ToLower("mmd_" + s.Alias)
}
