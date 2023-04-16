package main

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
)

func main() {
	siteURL := getSiteURL()
	strategy := getStrategy(siteURL)
	creds, err := getCreds(siteURL, strategy)
	if err != nil {
		fmt.Printf("\033[31mInvalid strategy: %s\033[0m\n", err)
		return
	}
	_, err = checkAuth(siteURL, strategy, creds)
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}
}

func getSiteURL() string {
	siteURLQ := &survey.Input{
		Message: "SharePoint URL:",
		Help:    "Site absolute URL, e.g. https://contoso.sharepoint.com/sites/MySite",
	}

	var siteURL string
	survey.AskOne(siteURLQ, &siteURL, survey.WithValidator(shouldBeURL))
	return siteURL
}

func getStrategy(siteURL string) string {
	strats, err := action("Resolving auth strategy...", func() ([]string, error) {
		return getStrategies(siteURL)
	})
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
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

	return strategy
}

func getCreds(siteURL, strategy string) ([][]string, error) {
	s, ok := stratsConf[strategy]
	if !ok {
		return nil, fmt.Errorf("can't resolve strategy %s", strategy)
	}
	return s.Creds(), nil
}

func checkAuth(siteURL, strategy string, creds [][]string) (*api.SP, error) {
	auth, err := newAuthByStrategy(strategy)
	if err != nil {
		return nil, err
	}

	cnfg := map[string]string{"siteURL": siteURL}
	for _, c := range creds {
		cnfg[c[0]] = c[1]
	}
	credsBytes, _ := json.Marshal(cnfg)

	if err := auth.ParseConfig(credsBytes); err != nil {
		return nil, err
	}

	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)

	web, err := action("Reaching site, checking auth...", sp.Web().Get)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\033[32mSuccess! Site title: \"%s\"\033[0m\n", web.Data().Title)

	return sp, nil
}

func action[T any](message string, fn func() (T, error)) (T, error) {
	fmt.Printf("\033[33m%s\033[0m", message)
	data, err := fn()
	if err != nil {
		fmt.Print("\033[2K\r")
		return data, err
	}
	fmt.Print("\033[2K\r")
	return data, nil
}
