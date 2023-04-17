package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip/cpass"
)

type credsSurv struct{}

var crypt = cpass.Cpass("")
var credsResolver = &credsSurv{}

func (*credsSurv) ondemand() [][]string {
	return [][]string{}
}

func (*credsSurv) user() [][]string {
	var username string
	interuptable(survey.AskOne(&survey.Input{
		Message: "User name:",
	}, &username, survey.WithValidator(survey.Required)))

	var password string
	interuptable(survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, survey.WithValidator(survey.Required)))

	var encrypt bool
	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}, &encrypt))

	if encrypt {
		password, _ = crypt.Encode(password)
	}

	return [][]string{
		{"username", username},
		{"password", password},
	}
}

func (*credsSurv) ntlm() [][]string {
	var domain string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Domain:",
	}, &domain, survey.WithValidator(survey.Required)))

	return append([][]string{{"domain", domain}}, credsResolver.user()...)
}

func (*credsSurv) saml() [][]string {
	var username string
	interuptable(survey.AskOne(&survey.Input{
		Message: "User name:",
	}, &username, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeEmail)))

	var password string
	interuptable(survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, survey.WithValidator(survey.Required)))

	var encrypt bool
	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}, &encrypt))

	if encrypt {
		password, _ = crypt.Encode(password)
	}

	return [][]string{
		{"username", username},
		{"password", password},
	}
}

func (*credsSurv) azurebase() [][]string {
	var tenantID string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Tenant ID:",
	}, &tenantID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID)))

	var clientID string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Client ID:",
	}, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID)))

	return [][]string{
		{"tenantId", tenantID},
		{"clientId", clientID},
	}
}

func (c *credsSurv) azurecert() [][]string {
	azurebase := c.azurebase()

	var certPath string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Certificate path:",
	}, &certPath, survey.WithValidator(survey.Required)))

	var certPass string
	interuptable(survey.AskOne(&survey.Password{
		Message: "Certificate password:",
	}, &certPass, survey.WithValidator(survey.Required)))

	var encrypt bool
	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}, &encrypt))

	if encrypt {
		certPass, _ = crypt.Encode(certPass)
	}

	return append(azurebase, [][]string{
		{"certPath", certPath},
		{"certPass", certPass},
	}...)
}

func (c *credsSurv) azurecreds() [][]string {
	return append(c.azurebase(), credsResolver.saml()...)
}

func (c *credsSurv) device() [][]string {
	return c.azurebase()
}

func (*credsSurv) addin() [][]string {
	var clientID string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Client ID:",
	}, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID)))

	var clientSecret string
	interuptable(survey.AskOne(&survey.Password{
		Message: "Client Secret:",
	}, &clientSecret, survey.WithValidator(survey.Required)))

	var encrypt bool
	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Encrypt secret?",
		Default: true,
	}, &encrypt))

	if encrypt {
		clientSecret, _ = crypt.Encode(clientSecret)
	}

	return [][]string{
		{"clientId", clientID},
		{"clientSecret", clientSecret},
	}
}

func (*credsSurv) adfs() [][]string {
	var username string
	interuptable(survey.AskOne(&survey.Input{
		Message: "User name:",
	}, &username, survey.WithValidator(survey.Required)))

	var password string
	interuptable(survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, survey.WithValidator(survey.Required)))

	var relyingParty string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Relying Party:",
		Default: "urn:sharepoint:www",
	}, &relyingParty, survey.WithValidator(survey.Required)))

	var adfsURL string
	interuptable(survey.AskOne(&survey.Input{
		Message: "ADFS URL:",
		Help:    "E.g.: https://login.contoso.com",
	}, &adfsURL, survey.WithValidator(survey.Required)))

	var adfsCookie string
	interuptable(survey.AskOne(&survey.Input{
		Message: "ADFS Cookie:",
		Default: "FedAuth",
	}, &adfsCookie, survey.WithValidator(survey.Required)))

	var encrypt bool
	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}, &encrypt))

	if encrypt {
		password, _ = crypt.Encode(password)
	}

	return [][]string{
		{"username", username},
		{"password", password},
		{"relyingParty", relyingParty},
		{"adfsUrl", adfsURL},
		{"adfsCookie", adfsCookie},
	}
}
