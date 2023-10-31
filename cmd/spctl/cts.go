package main

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/ct"
	"github.com/koltyakov/gosip/api"
)

type ContentTypeConf struct {
	ID   string
	Spec ct.Spec
}

func getContentTypesConf(sp *api.SP) ([]ContentTypeConf, error) {
	resp, err := action("Getting content types", func() (api.ContentTypesResp, error) {
		return sp.Web().ContentTypes().
			Top(5000).
			Filter("Hidden eq false and Group ne '_Hidden'").
			OrderBy("Id", true).
			Get()
	})
	if err != nil {
		return nil, err
	}

	data := resp.Data()
	ctt := make([]string, len(data))
	cttMap := map[string]api.ContentTypeResp{}
	for i, t := range data {
		c := t.Data()
		ctKey := c.Name + " \033[90m(" + c.Group + ") [" + c.ID + "]\033[0m"

		cttMap[c.ID] = t
		ctt[i] = ctKey
	}

	var contentTypesToSync []string
	interuptable(survey.AskOne(&survey.MultiSelect{
		Message: "Select content types to sync:",
		Options: ctt,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}, &contentTypesToSync, survey.WithValidator(survey.Required)))

	ctConf := make([]ContentTypeConf, len(contentTypesToSync))
	for i, t := range contentTypesToSync {
		fieldsConf, err := getFieldsConf(sp, t)
		if err != nil {
			return nil, err
		}
		ctInfo, ok := cttMap[getEntityID(t)]
		if !ok {
			return nil, err
		}
		ctConf[i] = ContentTypeConf{
			ID: ctInfo.Data().Name,
			Spec: ct.Spec{
				Select: fieldsConf.Select,
				Expand: fieldsConf.Expand,
			},
		}
	}

	return ctConf, nil
}
