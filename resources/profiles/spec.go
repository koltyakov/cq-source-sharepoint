package profiles

// Spec is the configuration for MMD term set source
type Spec struct {
	// Whether to enable user profiles sync
	Enabled bool `json:"enabled"`
	// Optional, an alias for the table name
	Alias string `json:"alias"`
}

// SetDefault sets default values for list spec
func (*Spec) SetDefault() {
	// Default values
}

// Validate validates user profiles spec validity
func (s *Spec) Validate() error {
	// Nothing to validate
	return nil
}

// GetAlias returns the alias for the term set
func (s *Spec) GetAlias() string {
	if s.Alias == "" {
		return "ups_profile"
	}
	return "ups_" + s.Alias
}
