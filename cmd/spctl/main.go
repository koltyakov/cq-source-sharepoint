package main

import (
	"fmt"
	"net/url"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	siteURLQ := &survey.Input{
		Message: "SharePoint URL:",
		Help:    "Site absolute URL, e.g. https://contoso.sharepoint.com/sites/MySite",
	}

	var siteURL string
	survey.AskOne(siteURLQ, &siteURL, survey.WithValidator(shouldBeURL))

	// Print in yellow
	fmt.Print("\033[33m" + "Resolving auth strategy..." + "\033[0m")

	strats, err := getStrategies(siteURL)
	if err != nil {
		fmt.Print("\033[2K\r") // Clear line
		fmt.Println("\033[31mError: " + err.Error() + "\033[0m")
		// fmt.Println("\033[31m" + "Failed to resolve auth strategy, using all available" + "\033[0m")
		strats = allStrats
	}

	strategyQ := &survey.Select{
		Message: "Authentication method:",
		Options: strats,
		Help:    "See more at https://go.spflow.com/auth/overview",
		Description: func(value string, index int) string {
			return stratsConf[value].Description
		},
	}

	var strategy string
	survey.AskOne(strategyQ, &strategy)
}

func shouldBeURL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}

	if _, err := url.ParseRequestURI(str); err != nil {
		return fmt.Errorf("value is not a valid URL")
	}

	return nil
}
