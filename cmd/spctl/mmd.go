package main

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/cq-source-sharepoint/resources/mmd"
	"github.com/koltyakov/gosip/api"
)

type MMDConf struct {
	ID   string
	Spec mmd.Spec
}

type termSetInfo struct {
	ID    string
	Name  string
	Group string
}

func getMMDConf(sp *api.SP) ([]MMDConf, error) {
	store := sp.Taxonomy().Stores().Default()
	groups, err := store.Groups().Get()
	if err != nil {
		return nil, err
	}

	termsetsNames := []string{}
	termSets := map[string]termSetInfo{}
	for _, group := range groups {
		sets, err := store.Groups().GetByID(extractMmdID(group["Id"].(string))).Sets().Get()
		if err != nil {
			return nil, err
		}
		for _, set := range sets {
			tsID := extractMmdID(set["Id"].(string))
			tsName := set["Name"].(string)
			tsGroup := group["Name"].(string)
			termSets[tsID] = termSetInfo{
				ID:    tsID,
				Name:  tsName,
				Group: tsGroup,
			}
			termsetsName := tsName + " \033[90m(" + tsGroup + ") [" + tsID + "]\033[0m"
			termsetsNames = append(termsetsNames, termsetsName)
		}
	}

	var termSetsToSync []string
	interuptable(survey.AskOne(&survey.MultiSelect{
		Message: "Select term sets to sync:",
		Options: termsetsNames,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}, &termSetsToSync, survey.WithValidator(survey.Required)))

	mmdConf := make([]MMDConf, len(termSetsToSync))
	for i, t := range termSetsToSync {
		termSetID := getEntityID(t)
		mmdConf[i] = MMDConf{
			ID: termSetID,
			Spec: mmd.Spec{
				Alias: termSets[termSetID].Name,
			},
		}
	}

	return mmdConf, nil
}

func extractMmdID(idString string) string {
	// /Guid(062962b3-82e0-4416-9b62-41efba8e23db)/ -> 062962b3-82e0-4416-9b62-41efba8e23db
	return idString[6 : len(idString)-2]
}
