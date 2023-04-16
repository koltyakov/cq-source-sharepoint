package main

type Strategy struct {
	Desc  string
	Docs  string
	Envs  []string
	Creds func() [][]string
}

const (
	SPO    string = "spo"
	OnPrem        = "onprem"
)

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

var stratsConf = map[string]Strategy{
	"ondemand": {
		Desc:  "On-Demand (Browser Popup Prompt)",
		Docs:  "https://go.spflow.com/auth/strategies/on-demand",
		Envs:  []string{SPO, OnPrem},
		Creds: creds.ondemand,
	},
	"azurecert": {
		Desc:  "Azure App (Certificate) [SPO]",
		Docs:  "https://go.spflow.com/auth/strategies/azure-certificate-auth",
		Envs:  []string{SPO},
		Creds: creds.azurecert,
	},
	"azurecreds": {
		Desc:  "Azure App (Client Credentials) [SPO]",
		Docs:  "https://go.spflow.com/auth/strategies/azure-creds-auth",
		Envs:  []string{SPO},
		Creds: creds.azurecreds,
	},
	"device": {
		Desc:  "Azure App (Device Login) [SPO]",
		Docs:  "https://go.spflow.com/auth/strategies/azure-device-flow",
		Envs:  []string{SPO},
		Creds: creds.device,
	},
	"addin": {
		Desc:  "Add-In Only (Legacy) [SPO]",
		Docs:  "https://go.spflow.com/auth/strategies/addin",
		Envs:  []string{SPO},
		Creds: creds.addin,
	},
	"adfs": {
		Desc:  "ADFS [SPO, On-Premises]",
		Docs:  "https://go.spflow.com/auth/strategies/adfs",
		Envs:  []string{SPO, OnPrem},
		Creds: creds.user,
	},
	"fba": {
		Desc:  "FBA (Legacy) [On-Premises]",
		Docs:  "https://go.spflow.com/auth/strategies/fba",
		Envs:  []string{OnPrem},
		Creds: creds.user,
	},
	"ntlm": {
		Desc:  "NTLM [On-Premises]",
		Docs:  "https://go.spflow.com/auth/strategies/ntlm",
		Envs:  []string{OnPrem},
		Creds: creds.ntlm,
	},
	"ntlm2": {
		Desc:  "NTLM (Alternative) [On-Premises]",
		Docs:  "https://go.spflow.com/auth/strategies/alternative-ntlm",
		Envs:  []string{OnPrem},
		Creds: creds.ntlm,
	},
	"saml": {
		Desc:  "SAML (Client Credentials) [SPO]",
		Docs:  "https://go.spflow.com/auth/strategies/saml",
		Envs:  []string{SPO},
		Creds: creds.saml,
	},
	"tmg": {
		Desc:  "TMG (Legacy) [On-Premises]",
		Docs:  "https://go.spflow.com/auth/strategies/tmg",
		Envs:  []string{OnPrem},
		Creds: creds.user,
	},
}
