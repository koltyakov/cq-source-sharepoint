package auth

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	"github.com/koltyakov/gosip/auth"
)

func GetSP(spec Spec) (*api.SP, error) {
	jsonCreds, _ := json.Marshal(spec.Creds)
	authCnfg, err := auth.NewAuthByStrategy(spec.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth config: %w", err)
	}
	if err := authCnfg.ParseConfig(jsonCreds); err != nil {
		return nil, fmt.Errorf("failed to parse auth config: %w", err)
	}

	return api.NewSP(&gosip.SPClient{AuthCnfg: authCnfg}), nil
}
