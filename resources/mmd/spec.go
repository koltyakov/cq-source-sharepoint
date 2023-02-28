package mmd

// Spec is the configuration for MMD term set source
type Spec struct {
	// Optional, an alias for the table name
	// Don't map different lists to the same table - such scenariou is not supported
	Alias string `json:"alias"`
}

// SetDefault sets default values for list spec
// func (s *Spec) SetDefault() {}
