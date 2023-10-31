package main

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/lists"
	"github.com/koltyakov/gosip/api"
)

type ListConf struct {
	ID   string
	Spec lists.Spec
}

type listInfo struct {
	ID         string `json:"Id"`
	Title      string `json:"Title"`
	RootFolder struct {
		URL string `json:"ServerRelativeUrl"`
	} `json:"RootFolder"`
}

func getListsConf(sp *api.SP) ([]ListConf, error) {
	resp, err := action("Getting lists", func() (api.ListsResp, error) {
		return sp.Web().Lists().
			Top(5000).
			Select("Id,Title,RootFolder/ServerRelativeUrl").
			Expand("RootFolder").Get()
	})
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(sp.ToURL())
	basePath := u.Path + "/"

	data := resp.Data()
	ll := make([]string, len(data))
	llMap := map[string]listInfo{}
	for i, l := range data {
		info := listInfo{}
		_ = json.Unmarshal(l.Normalized(), &info)
		listURI := strings.Replace(info.RootFolder.URL, basePath, "", 1)

		listKey := info.Title + " \033[90m[" + listURI + "]\033[0m"
		llMap[listURI] = info
		ll[i] = listKey
	}

	var listsToSync []string
	interuptable(survey.AskOne(&survey.MultiSelect{
		Message: "Select lists to sync:",
		Options: ll,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}, &listsToSync, survey.WithValidator(survey.Required)))

	listsConf := make([]ListConf, len(listsToSync))
	for i, l := range listsToSync {
		listURI := getEntityID(l)
		fieldsConf, err := getFieldsConf(sp, l)
		if err != nil {
			return nil, err
		}
		listsConf[i] = ListConf{
			ID: listURI,
			Spec: lists.Spec{
				Select: fieldsConf.Select,
				Expand: fieldsConf.Expand,
			},
		}
	}

	return listsConf, nil
}
