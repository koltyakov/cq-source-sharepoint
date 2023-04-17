package main

import "github.com/AlecAivazis/survey/v2"

var syncScenariosMap = map[string]string{
	"Lists and libraries":    "lists",
	"Content types rollup":   "content_types",
	"Managed metadata terms": "mmd",
	"User profiles (UPS)":    "profiles",
	"Search driven data":     "search",
}

func getSyncScenarios() []string {
	var syncScenarios []string
	interuptable(survey.AskOne(&survey.MultiSelect{
		Message: "Select subjects of sync:",
		Options: []string{
			"Lists and libraries",
			"Content types rollup",
			"Managed metadata terms",
			"User profiles (UPS)",
			"Search driven data",
		},
	}, &syncScenarios, survey.WithValidator(survey.Required)))

	for i, s := range syncScenarios {
		syncScenarios[i] = syncScenariosMap[s]
	}

	return syncScenarios
}
