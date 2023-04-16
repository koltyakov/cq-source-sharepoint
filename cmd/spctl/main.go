package main

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
)

func main() {
	siteURLQ := &survey.Input{
		Message: "SharePoint URL:",
		Help:    "Site absolute URL, e.g. https://contoso.sharepoint.com/sites/MySite",
	}

	var siteURL string
	survey.AskOne(siteURLQ, &siteURL, survey.WithValidator(shouldBeURL))

	fmt.Print("\033[33m" + "Resolving auth strategy..." + "\033[0m")

	strats, err := getStrategies(siteURL)
	if err != nil {
		fmt.Print("\033[2K\r") // Clear line
		fmt.Println("\033[31mError: " + err.Error() + "\033[0m")
		// fmt.Println("\033[31m" + "Failed to resolve auth strategy, using all available" + "\033[0m")
		strats = allStrats
	}

	strategyQ := &survey.Select{
		Message: "Auth method:",
		Options: strats,
		Help:    "See more at https://go.spflow.com/auth/overview",
		Description: func(value string, index int) string {
			return stratsConf[value].Desc
		},
	}

	var strategy string
	survey.AskOne(strategyQ, &strategy)

	s, ok := stratsConf[strategy]
	if !ok {
		fmt.Println("\033[31m" + "Invalid strategy" + "\033[0m")
		return
	}

	auth, err := newAuthByStrategy(strategy)
	if err != nil {
		fmt.Println("\033[31m" + "Error: " + err.Error() + "\033[0m")
		return
	}

	authCreds := s.Creds()

	credsConfig := map[string]string{
		"siteURL": siteURL,
	}
	for _, c := range authCreds {
		credsConfig[c[0]] = c[1]
	}
	credsBytes, _ := json.Marshal(credsConfig)

	if err := auth.ParseConfig(credsBytes); err != nil {
		fmt.Println("\033[31m" + "Error: " + err.Error() + "\033[0m")
		return
	}

	fmt.Print("\033[33m" + "Reaching site, checking auth..." + "\033[0m")

	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)

	web, err := sp.Web().Get()
	if err != nil {
		fmt.Print("\033[2K\r") // Clear line
		fmt.Println("\033[31mError: " + err.Error() + "\033[0m")
		return
	}

	fmt.Print("\033[2K\r") // Clear line
	fmt.Println("\033[32m" + "Success! Site title: \"" + web.Data().Title + "\"\033[0m")
}
