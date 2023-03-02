package search

// Spec is the configuration for Search source
type Spec struct {
	QueryText        string   `json:"query_text"`        // Required
	TrimDuplicates   bool     `json:"trim_duplicates"`   // Optional, default is false
	SourceID         string   `json:"source_id"`         // Optional, default is empty
	SelectProperties []string `json:"select_properties"` // Optional, default is empty array
}

// SetDefault sets default values for list spec
func (s *Spec) SetDefault() {
	if s.QueryText == "" {
		s.QueryText = "*"
	}
}
