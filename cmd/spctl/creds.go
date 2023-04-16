package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip/cpass"
)

type credsSurv struct{}

var crypt = cpass.Cpass("")
var creds = &credsSurv{}

func (c *credsSurv) ondemand() [][]string {
	return [][]string{}
}

func (c *credsSurv) user() [][]string {
	var username string
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	survey.AskOne(encryptQ, &encrypt)

	if encrypt {
		password, _ = crypt.Encode(password)
	}

	return [][]string{
		{"username", username},
		{"password", password},
	}
}

func (c *credsSurv) ntlm() [][]string {
	var domain string
	domainQ := &survey.Input{
		Message: "Domain:",
	}
	survey.AskOne(domainQ, &domain, survey.WithValidator(survey.Required))

	return append([][]string{{"domain", domain}}, creds.user()...)
}

func (c *credsSurv) saml() [][]string {
	var username string
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeEmail))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	survey.AskOne(encryptQ, &encrypt)

	if encrypt {
		password, _ = crypt.Encode(password)
	}

	return [][]string{
		{"username", username},
		{"password", password},
	}
}

func (c *credsSurv) azurebase() [][]string {
	var tenantID string
	tenantIDQ := &survey.Input{
		Message: "Tenant ID:",
	}
	survey.AskOne(tenantIDQ, &tenantID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

	var clientID string
	clientIDQ := &survey.Input{
		Message: "Client ID:",
	}
	survey.AskOne(clientIDQ, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

	return [][]string{
		{"tenantId", tenantID},
		{"clientId", clientID},
	}
}

func (c *credsSurv) azurecert() [][]string {
	azurebase := c.azurebase()

	var certPath string
	certPathQ := &survey.Input{
		Message: "Certificate path:",
	}
	survey.AskOne(certPathQ, &certPath, survey.WithValidator(survey.Required))

	var certPass string
	certPassQ := &survey.Password{
		Message: "Certificate password:",
	}
	survey.AskOne(certPassQ, &certPass, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	survey.AskOne(encryptQ, &encrypt)

	if encrypt {
		certPass, _ = crypt.Encode(certPass)
	}

	return append(azurebase, [][]string{
		{"certPath", certPath},
		{"certPass", certPass},
	}...)
}

func (c *credsSurv) azurecreds() [][]string {
	return append(c.azurebase(), creds.saml()...)
}

func (c *credsSurv) device() [][]string {
	return c.azurebase()
}

func (c *credsSurv) addin() [][]string {
	var clientID string
	clientIDQ := &survey.Input{
		Message: "Client ID:",
	}
	survey.AskOne(clientIDQ, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

	var clientSecret string
	clientSecretQ := &survey.Password{
		Message: "Client Secret:",
	}
	survey.AskOne(clientSecretQ, &clientSecret, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt secret?",
		Default: true,
	}
	survey.AskOne(encryptQ, &encrypt)

	if encrypt {
		clientSecret, _ = crypt.Encode(clientSecret)
	}

	return [][]string{
		{"clientId", clientID},
		{"clientSecret", clientSecret},
	}
}

func (c *credsSurv) adfs() [][]string {
	var username string
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var relyingParty string
	relyingPartyQ := &survey.Input{
		Message: "Relying Party:",
		Default: "urn:sharepoint:www",
	}
	survey.AskOne(relyingPartyQ, &relyingParty, survey.WithValidator(survey.Required))

	var adfsURL string
	adfsURLQ := &survey.Input{
		Message: "ADFS URL:",
		Help:    "E.g.: https://login.contoso.com",
	}
	survey.AskOne(adfsURLQ, &adfsURL, survey.WithValidator(survey.Required))

	var adfsCookie string
	adfsCookieQ := &survey.Input{
		Message: "ADFS Cookie:",
		Default: "FedAuth",
	}
	survey.AskOne(adfsCookieQ, &adfsCookie, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	survey.AskOne(encryptQ, &encrypt)

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
