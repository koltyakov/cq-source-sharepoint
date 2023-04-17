package main

import "github.com/AlecAivazis/survey/v2"

var syncScenariosMap = map[string]string{
	"Lists and libraries":    "lists",
	"Content types rollup":   "content_types",
	"Search driven queries":  "search",
	"Managed metadata terms": "mmd",
	"User profiles (UPS)":    "profiles",
}

func getSyncScenarios() []string {
	syncScenariosQ := &survey.MultiSelect{
		Message: "Select subjects of sync:",
		Options: []string{
			"Lists and libraries",
			"Content types rollup",
			"Search driven queries",
			"Managed metadata terms",
			"User profiles (UPS)",
		},
	}

	var syncScenarios []string
	_ = survey.AskOne(syncScenariosQ, &syncScenarios, survey.WithValidator(survey.Required))

	for i, s := range syncScenarios {
		syncScenarios[i] = syncScenariosMap[s]
	}

	return syncScenarios
}
