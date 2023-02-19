package client

import "fmt"

type Spec struct {
	SiteURL 		 string `json:"site_url"`
	ClientID 		 string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Lists 			 []string `json:"lists"`
	// plugin spec goes here
}


func (s *Spec) Validate() error {
	if s.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if s.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	if s.SiteURL == "" {
		return fmt.Errorf("site_url is required")
	}
	return nil
}