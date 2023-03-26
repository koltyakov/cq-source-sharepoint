package profiles

// Spec is the configuration for MMD term set source
type Spec struct {
	// Whether to enable user profiles sync
	Enabled bool `json:"enabled"`
	// Optional, an alias for the table name
	Alias string `json:"alias"`
}

// SetDefault sets default values for list spec
func (*Spec) SetDefault() {}
