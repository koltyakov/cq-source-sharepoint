package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
)

var pluginVersion = "v1.6.2"

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		panic("interrupted")
	}()
}

func main() {
	siteURL := getSiteURL()
	strategy := getStrategy(siteURL)
	creds, err := getCreds(strategy)
	if err != nil {
		fmt.Printf("\033[31mInvalid strategy: %s\033[0m\n", err)
		return
	}
	sp, err := checkAuth(siteURL, strategy, creds)
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}

	version, _ := getPluginVersion()
	source := getSourceName()
	destination := getDestination()

	spec := &SourceSpec{
		Name:         source,
		Registry:     "github",
		Path:         "koltyakov/sharepoint",
		Version:      version,
		Destinations: []string{destination},
		Spec: PluginSpec{
			Auth: AuthSpec{
				Strategy: strategy,
				Creds:    append([][]string{{"siteUrl", siteURL}}, creds...),
			},
		},
	}

	syncScenarios := getSyncScenarios()
	for _, scenario := range syncScenarios {
		if scenario == "lists" {
			listsConf, err := getListsConf(sp)
			if err != nil {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}
			spec.Spec.Lists = listsConf
		}

		if scenario == "content_types" {
			contentTypesConf, err := getContentTypesConf(sp)
			if err != nil {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}
			spec.Spec.ContentTypes = contentTypesConf
		}

		if scenario == "mmd" {
			mmdConf, err := getMMDConf(sp)
			if err != nil {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}
			spec.Spec.MMD = mmdConf
		}

		if scenario == "search" {
			searchConfs, err := getSearchConfs(sp)
			if err != nil {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}
			spec.Spec.Search = searchConfs
		}

		if scenario == "profiles" {
			spec.Spec.Profiles = true
		}
	}

	if err := spec.Save(source + ".yml"); err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}
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

func getSiteURL() string {
	var siteURL string
	interuptable(survey.AskOne(&survey.Input{
		Message: "SharePoint URL:",
		Help:    "Site absolute URL, e.g. https://contoso.sharepoint.com/sites/MySite",
	}, &siteURL, survey.WithValidator(shouldBeURL), survey.WithValidator(shouldBeSPSite)))
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

	var strategy string
	interuptable(survey.AskOne(&survey.Select{
		Message: "Auth method:",
		Options: strats,
		Help:    "See more at https://go.spflow.com/auth/overview",
		Description: func(value string, index int) string {
			return stratsConf[value].Desc
		},
	}, &strategy))

	return strategy
}

func getCreds(strategy string) ([][]string, error) {
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

	if web.Data().Title == "" {
		return nil, fmt.Errorf("can't reach site, check that Site URL is correct")
	}

	fmt.Printf("\033[32mSuccess! Site title: \"%s\"\033[0m\n", web.Data().Title)

	return sp, nil
}

func getSourceName() string {
	var sourceName string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Source name:",
		Default: "sharepoint",
		Help:    "Source name to be used in the config file",
	}, &sourceName, survey.WithValidator(survey.Required)))
	return sourceName
}

func getDestination() string {
	var destination string
	interuptable(survey.AskOne(&survey.Input{
		Message: "Destination name:",
		Default: "postgres",
		Help:    "Destination name to be used in the config file",
	}, &destination, survey.WithValidator(survey.Required)))
	return destination
}
