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
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	_ = survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	_ = survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	_ = survey.AskOne(encryptQ, &encrypt)

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
	domainQ := &survey.Input{
		Message: "Domain:",
	}
	_ = survey.AskOne(domainQ, &domain, survey.WithValidator(survey.Required))

	return append([][]string{{"domain", domain}}, credsResolver.user()...)
}

func (*credsSurv) saml() [][]string {
	var username string
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	_ = survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeEmail))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	_ = survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	_ = survey.AskOne(encryptQ, &encrypt)

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
	tenantIDQ := &survey.Input{
		Message: "Tenant ID:",
	}
	_ = survey.AskOne(tenantIDQ, &tenantID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

	var clientID string
	clientIDQ := &survey.Input{
		Message: "Client ID:",
	}
	_ = survey.AskOne(clientIDQ, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

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
	_ = survey.AskOne(certPathQ, &certPath, survey.WithValidator(survey.Required))

	var certPass string
	certPassQ := &survey.Password{
		Message: "Certificate password:",
	}
	_ = survey.AskOne(certPassQ, &certPass, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	_ = survey.AskOne(encryptQ, &encrypt)

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
	clientIDQ := &survey.Input{
		Message: "Client ID:",
	}
	_ = survey.AskOne(clientIDQ, &clientID, survey.WithValidator(survey.Required), survey.WithValidator(shouldBeGUID))

	var clientSecret string
	clientSecretQ := &survey.Password{
		Message: "Client Secret:",
	}
	_ = survey.AskOne(clientSecretQ, &clientSecret, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt secret?",
		Default: true,
	}
	_ = survey.AskOne(encryptQ, &encrypt)

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
	usernameQ := &survey.Input{
		Message: "User name:",
	}
	_ = survey.AskOne(usernameQ, &username, survey.WithValidator(survey.Required))

	var password string
	passwordQ := &survey.Password{
		Message: "Password:",
	}
	_ = survey.AskOne(passwordQ, &password, survey.WithValidator(survey.Required))

	var relyingParty string
	relyingPartyQ := &survey.Input{
		Message: "Relying Party:",
		Default: "urn:sharepoint:www",
	}
	_ = survey.AskOne(relyingPartyQ, &relyingParty, survey.WithValidator(survey.Required))

	var adfsURL string
	adfsURLQ := &survey.Input{
		Message: "ADFS URL:",
		Help:    "E.g.: https://login.contoso.com",
	}
	_ = survey.AskOne(adfsURLQ, &adfsURL, survey.WithValidator(survey.Required))

	var adfsCookie string
	adfsCookieQ := &survey.Input{
		Message: "ADFS Cookie:",
		Default: "FedAuth",
	}
	_ = survey.AskOne(adfsCookieQ, &adfsCookie, survey.WithValidator(survey.Required))

	var encrypt bool
	encryptQ := &survey.Confirm{
		Message: "Encrypt password?",
		Default: true,
	}
	_ = survey.AskOne(encryptQ, &encrypt)

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
