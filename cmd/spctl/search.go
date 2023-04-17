package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip/api"
)

type SearchConf struct {
	ID               string
	QueryText        string
	SourceID         string
	TrimDuplicates   bool
	SelectProperties []string
}

func getSearchConfs(sp *api.SP) ([]SearchConf, error) {
	searchConfs := []SearchConf{}
	more := true

	for more {
		searchConf, err := getSearchConf(sp)
		if err != nil {
			return searchConfs, err
		}
		searchConfs = append(searchConfs, searchConf)

		interuptable(survey.AskOne(&survey.Confirm{
			Message: "Add another search?",
			Default: false,
		}, &more))
	}

	return searchConfs, nil
}

func getSearchConf(sp *api.SP) (SearchConf, error) {
	searchConf := SearchConf{}

	interuptable(survey.AskOne(&survey.Input{
		Message: "Search name:",
		Help:    "Search query name, e.g. 'Documents', used as table alias",
	}, &searchConf.ID, survey.WithValidator(survey.Required)))

	interuptable(survey.AskOne(&survey.Input{
		Message: "Query text:",
		Default: "*",
		Help:    "Search query text, see more https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference",
	}, &searchConf.QueryText, survey.WithValidator(survey.Required)))

	interuptable(survey.AskOne(&survey.Input{
		Message: "Source ID (optional):",
		Help:    "Search source ID, see more https://learn.microsoft.com/en-us/sharepoint/manage-result-sources",
	}, &searchConf.SourceID, survey.WithValidator(shouldBeGUIDorEmpty)))

	interuptable(survey.AskOne(&survey.Confirm{
		Message: "Trim duplicates?",
		Default: true,
		Help:    "Trim duplicates from search results",
	}, &searchConf.TrimDuplicates))

	res, err := sp.Search().PostQuery(&api.SearchQuery{
		QueryText:      searchConf.QueryText,
		SourceID:       searchConf.SourceID,
		TrimDuplicates: searchConf.TrimDuplicates,
		RowLimit:       1,
	})
	if err != nil {
		return searchConf, err
	}

	rows := res.Data().PrimaryQueryResult.RelevantResults.Table.Rows
	if len(rows) == 0 {
		fmt.Println("Warning: no results found for the query")
		return searchConf, nil
	}

	searchProps := []string{}
	for _, prop := range rows[0].Cells {
		propName := prop.Key
		searchProps = append(searchProps, propName)
	}

	interuptable(survey.AskOne(&survey.MultiSelect{
		Message: "Select properties:",
		Options: searchProps,
		// Default: "Title,Path,Author,Editor,LastModifiedTime,Size",
	}, &searchConf.SelectProperties, survey.WithValidator(survey.Required)))

	return searchConf, nil
}
