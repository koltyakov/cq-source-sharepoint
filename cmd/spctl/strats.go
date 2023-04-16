package main

type Strategy struct {
	Description string
	Docs        string
	Envs        []string
}

const (
	SPO    string = "spo"
	OnPrem        = "onprem"
)

var stratsConf = map[string]Strategy{
	"ondemand": {
		Description: "On-Demand (Browser Popup Prompt)",
		Docs:        "https://go.spflow.com/auth/strategies/on-demand",
		Envs:        []string{SPO, OnPrem},
	},
	"azurecert": {
		Description: "Azure App (Certificate) [SPO]",
		Docs:        "https://go.spflow.com/auth/strategies/azure-certificate-auth",
		Envs:        []string{SPO},
	},
	"azurecreds": {
		Description: "Azure App (Client Credentials) [SPO]",
		Docs:        "https://go.spflow.com/auth/strategies/azure-creds-auth",
		Envs:        []string{SPO},
	},
	"device": {
		Description: "Azure App (Device Login) [SPO]",
		Docs:        "https://go.spflow.com/auth/strategies/azure-device-flow",
		Envs:        []string{SPO},
	},
	"addin": {
		Description: "Add-In Only (Legacy) [SPO]",
		Docs:        "https://go.spflow.com/auth/strategies/addin",
		Envs:        []string{SPO},
	},
	"adfs": {
		Description: "ADFS [SPO, On-Premises]",
		Docs:        "https://go.spflow.com/auth/strategies/adfs",
		Envs:        []string{SPO, OnPrem},
	},
	"fba": {
		Description: "FBA (Legacy) [On-Premises]",
		Docs:        "https://go.spflow.com/auth/strategies/fba",
		Envs:        []string{OnPrem},
	},
	"ntlm": {
		Description: "NTLM [On-Premises]",
		Docs:        "https://go.spflow.com/auth/strategies/ntlm",
		Envs:        []string{OnPrem},
	},
	"ntlm2": {
		Description: "NTLM (Alternative) [On-Premises]",
		Docs:        "https://go.spflow.com/auth/strategies/alternative-ntlm",
		Envs:        []string{OnPrem},
	},
	"saml": {
		Description: "SAML (Client Credentials) [SPO]",
		Docs:        "https://go.spflow.com/auth/strategies/saml",
		Envs:        []string{SPO},
	},
	"tmg": {
		Description: "TMG (Legacy) [On-Premises]",
		Docs:        "https://go.spflow.com/auth/strategies/tmg",
		Envs:        []string{OnPrem},
	},
}

var allStrats = []string{
	"ondemand",
	"azurecert",
	"azurecreds",
	"device",
	"saml",
	"addin",
	"adfs",
	"ntlm",
	"ntlm2",
	"fba",
	"tmg",
}
