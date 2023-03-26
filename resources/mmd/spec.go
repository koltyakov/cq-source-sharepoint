package mmd

// Spec is the configuration for MMD term set source
type Spec struct {
	// Optional, an alias for the table name
	// Don't map different term sets to the same table - such scenario is not supported
	Alias string `json:"alias"`
}

// SetDefault sets default values for MMD spec
func (s *Spec) SetDefault() {}
