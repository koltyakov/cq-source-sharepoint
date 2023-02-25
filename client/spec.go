package client

import (
	"fmt"

	"github.com/thoas/go-funk"
)

// Can we use camleCase here for JSON and snake for YML?
type Spec struct {
	// Gosip auth config connection params https://go.spflow.com/auth/overview
	Auth struct {
		Strategy string            `json:"strategy"`
		Creds    map[string]string `json:"creds"`
	} `json:"auth"`

	// Lists to fetch, if empty all lists will be fetched
	Lists map[string]ListSpec `json:"lists"`
}

type ListSpec struct {
	Select []string `json:"select"`
	Expand []string `json:"expand"`
	Alias  string   `json:"alias"`
}

func (s *Spec) SetDefaults() {
	if s.Lists == nil {
		s.Lists = make(map[string]ListSpec)
	}

	for ListURI, listSpec := range s.Lists {
		listSpec.SetDefault()
		s.Lists[ListURI] = listSpec
	}
}

// ToDo: Refactor to more elegant solution
func (l *ListSpec) SetDefault() {
	if l.Select == nil {
		l.Select = []string{}
	}

	exclude := []string{"*"}
	prepProps := []string{"ID"}
	apndProps := []string{"Created", "AuthorId", "Modified", "EditorId"}

	l.Select = funk.FilterString(l.Select, func(field string) bool {
		return !funk.ContainsString(concatSlice(exclude, concatSlice(prepProps, apndProps)), field)
	})

	l.Select = concatSlice(prepProps, concatSlice(l.Select, apndProps))
}

func (s Spec) Validate() error {
	if s.Auth.Strategy == "" {
		return fmt.Errorf("auth.strategy is required")
	}
	if len(s.Auth.Creds) == 0 {
		return fmt.Errorf("auth.creds is required")
	}

	return nil
}
