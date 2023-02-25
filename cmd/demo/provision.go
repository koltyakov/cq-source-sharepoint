package main

import (
	"encoding/xml"
	"strings"

	"github.com/koltyakov/gosip/api"
)

// EnsureList ensures list existence, returns true if list was created
func EnsureList(sp *api.SP, listTitle string, listURI string) (bool, error) {
	if listURI == "" {
		listURI = listTitle
	}

	if _, err := sp.Web().Lists().GetByTitle(listTitle).Get(); err != nil {
		if _, err := sp.Web().Lists().AddWithURI(listTitle, listURI, nil); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// EnsureListField ensures list field existence, returns true if field was created
func EnsureListField(sp *api.SP, listName string, fieldName string, schemaXML string) (bool, error) {
	fields := sp.Web().Lists().GetByTitle(listName).Fields()

	// Dynamically inject schema properties
	var schema struct {
		XMLName     xml.Name `xml:"Field"`
		DisplayName string   `xml:"DisplayName,attr"`
		Name        string   `xml:"Name,attr"`
		Type        string   `xml:"Type,attr"`
		List        string   `xml:"List,attr"`
	}
	if err := xml.Unmarshal([]byte(schemaXML), &schema); err == nil {
		// Autoinject field name
		if schema.Name == "" {
			schemaXML = strings.Replace(schemaXML, "<Field ", "<Field Name=\""+fieldName+"\" ", 1)
		}

		// Remplace display name with internal name for initial creation
		if schema.DisplayName != "" {
			schemaXML = strings.Replace(schemaXML, " DisplayName=\""+schema.DisplayName+"\"", " DisplayName=\""+fieldName+"\"", 1)
		}

		switch schema.Type {
		case "Lookup":
			// Resolve list by URI and inject list ID
			if list, err := sp.Web().GetList(schema.List).Select("Id").Get(); err == nil {
				schemaXML = strings.Replace(schemaXML, "List=\""+schema.List+"\"", "List=\""+list.Data().ID+"\"", 1)
			}
		}
	}

	if _, err := fields.GetByInternalNameOrTitle(fieldName).Get(); err != nil {
		if _, err := fields.CreateFieldAsXML(schemaXML, 16); err != nil {
			return false, err
		}
		if schema.DisplayName != "" {
			if _, err := fields.GetByInternalNameOrTitle(fieldName).Update([]byte(`{ "Title": "` + schema.DisplayName + `" }`)); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	return false, nil
}

// DropList deletes a list by title
func DropList(sp *api.SP, listTitle string) error {
	if err := sp.Web().Lists().GetByTitle(listTitle).Delete(); err != nil {
		return err
	}
	return nil
}
