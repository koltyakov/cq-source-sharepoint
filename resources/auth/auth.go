package auth

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
	ntlm2 "github.com/koltyakov/gosip-sandbox/strategies/ntlm"
	ondemand "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/api"
	"github.com/koltyakov/gosip/auth"
)

func GetSP(spec Spec) (*api.SP, error) {
	jsonCreds, _ := json.Marshal(spec.Creds)
	authCnfg, err := newAuthByStrategy(spec.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth config: %w", err)
	}
	if err := authCnfg.ParseConfig(jsonCreds); err != nil {
		return nil, fmt.Errorf("failed to parse auth config: %w", err)
	}

	return api.NewSP(&gosip.SPClient{AuthCnfg: authCnfg}), nil
}

func newAuthByStrategy(strategy string) (gosip.AuthCnfg, error) {
	// Some NTLM configuratios will need this auth instead of the default one
	if strategy == "ntlm2" {
		authCnfg := &ntlm2.AuthCnfg{}
		return authCnfg, nil
	}
	// Browser popup auth window (needs Chrome installed)
	if strategy == "ondemand" {
		authCnfg := &ondemand.AuthCnfg{}
		return authCnfg, nil
	}
	return auth.NewAuthByStrategy(strategy)
}
