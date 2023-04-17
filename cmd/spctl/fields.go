package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/gosip/api"
)

type FieldsConf struct {
	Select []string
	Expand []string
}

func getFieldsConf(sp *api.SP, entityName string) (*FieldsConf, error) {
	entityID := getEntityID(entityName)

	data, err := action("Getting fields for "+entityName, func() ([]*api.FieldInfo, error) {
		if getEntityType(entityID) == "list" {
			return getListFieldInfo(sp, entityID)
		}
		if getEntityType(entityID) == "content_type" {
			return getContentTypeFieldInfo(sp, entityID)
		}
		return nil, fmt.Errorf("unknown entity type: %s", getEntityType(entityID))
	})
	if err != nil {
		return nil, err
	}

	ignoreFields := []string{"AppAuthor", "AppEditor"}
	dd := []*api.FieldInfo{}
	for _, f := range data {
		if f.TypeAsString == "Lookup" && f.LookupList == "" {
			continue
		}

		if includes(ignoreFields, f.EntityPropertyName) {
			continue
		}

		dd = append(dd, f)
	}

	fieldsOptsStr := make([]string, len(dd))
	fieldsInfo := map[string]*api.FieldInfo{}
	for i, f := range dd {
		fieldsOptsStr[i] = f.Title +
			" \033[90m[" + f.EntityPropertyName + "]" +
			" " + f.TypeAsString + "\033[0m"
		fieldsInfo[f.EntityPropertyName] = f
	}

	defaultFieldNames := []string{"ID", "Title"}
	defaultFields := []string{}
	for _, f := range fieldsOptsStr {
		if includes(defaultFieldNames, getEntityID(f)) {
			defaultFields = append(defaultFields, f)
		}
	}

	// ToDo: No fields no prompt

	var fieldsToSync []string
	fieldsQ := &survey.MultiSelect{
		Message: "Select fields to sync for " + entityName + ":",
		Options: fieldsOptsStr,
		Default: defaultFields,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}
	_ = survey.AskOne(fieldsQ, &fieldsToSync, survey.WithValidator(survey.Required))

	fieldsConf := &FieldsConf{
		Select: make([]string, len(fieldsToSync)),
		Expand: []string{},
	}
	for i, f := range fieldsToSync {
		name := getEntityID(f)
		info, ok := fieldsInfo[name]
		if !ok {
			return nil, fmt.Errorf("field not found: %s", name)
		}
		fieldsConf.Select[i] = name
		if info.TypeAsString == "Lookup" || info.TypeAsString == "User" {
			fieldsConf.Select[i] = name + "/Id"
			fieldsConf.Expand = append(fieldsConf.Expand, name)
		}
	}

	return fieldsConf, nil
}

func getListFieldInfo(sp *api.SP, listURI string) ([]*api.FieldInfo, error) {
	resp, err := sp.Web().GetList(listURI).
		Fields().
		Filter("Hidden eq false and FieldTypeKind ne 12").
		Top(5000).
		Get()
	if err != nil {
		return nil, err
	}
	rr := resp.Data()
	dd := make([]*api.FieldInfo, len(rr))
	for i, r := range rr {
		dd[i] = r.Data()
	}
	return dd, nil
}

func getContentTypeFieldInfo(sp *api.SP, ctID string) ([]*api.FieldInfo, error) {
	type contentTypeInfo struct {
		Fields []*api.FieldInfo `json:"Fields"`
	}
	rest, err := sp.Web().ContentTypes().GetByID(ctID).
		Expand("Fields").
		Get()
	if err != nil {
		return nil, err
	}
	info := contentTypeInfo{}
	_ = json.Unmarshal(rest.Normalized(), &info)

	fields := []*api.FieldInfo{}
	for _, f := range info.Fields {
		if f.Hidden || f.FieldTypeKind == 12 {
			continue
		}
		fields = append(fields, f)
	}

	return fields, nil
}
