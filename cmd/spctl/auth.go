package main

import (
	"github.com/koltyakov/gosip"
	ntlm2 "github.com/koltyakov/gosip-sandbox/strategies/ntlm"
	"github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/auth"
)

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
