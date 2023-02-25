package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type listInfo struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	RootFolder  struct {
		ServerRelativeURL string `json:"ServerRelativeUrl"`
	} `json:"RootFolder"`
}

func (c *Client) getListInfo(listURI string) (*listInfo, error) {
	list := c.SP.Web().GetList(listURI)

	listResp, err := list.Select("Title,Description,RootFolder/ServerRelativeUrl").Expand("RootFolder").Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	var listInfo *listInfo
	if err := json.Unmarshal(listResp.Normalized(), &listInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list: %w", err)
	}

	return listInfo, nil
}

func (c *Client) tableFromList(listURI string, spec ListSpec) (*schema.Table, *ListModel, error) {
	listInfo, err := c.getListInfo(listURI)
	if err != nil {
		return nil, nil, err
	}

	tableName := normalizeName(listInfo.RootFolder.ServerRelativeURL)
	if spec.Alias != "" {
		tableName = normalizeName(spec.Alias)
	}

	table := &schema.Table{
		Name:        "sharepoint_" + tableName,
		Description: listInfo.Description,
	}
	logger := c.Logger.With().Str("table", table.Name).Logger()

	fields, err := c.SP.Web().GetList(listURI).Fields().Get()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get fields: %w", err)
	}

	fieldsData := fields.Data()
	mapping := &ListModel{
		ListURI:   listURI,
		ListSpec:  spec,
		FieldsMap: map[string]string{},
	}

	// ToDo: Rearchitect table construction logic
	for _, prop := range spec.Select {
		var field *api.FieldInfo
		for _, fieldResp := range fieldsData {
			fieldData := fieldResp.Data()
			propName := fieldData.EntityPropertyName
			lookups := []string{"Lookup", "User", "LookupMulti", "UserMulti"}
			if funk.Contains(lookups, fieldData.TypeAsString) {
				propName += "Id"
			}
			if propName == prop {
				field = fieldData
				break
			}
		}

		// Props is not presented in list's fields
		if field == nil {
			c := schema.Column{
				Name:        normalizeName(prop),
				Description: prop,
				Type:        schema.TypeString,
			}

			table.Columns = append(table.Columns, c)
			mapping.FieldsMap[c.Name] = prop
			continue
		}

		col := columnFromField(field, logger)
		col.CreationOptions.PrimaryKey = prop == "ID" // ToDo: Decide on ID cunstruction logic: use ID/UniqueID/Path+ID

		table.Columns = append(table.Columns, col)
		mapping.FieldsMap[col.Name] = prop
	}

	return table, mapping, nil
}

func columnFromField(field *api.FieldInfo, logger zerolog.Logger) schema.Column {
	c := schema.Column{
		Description: field.Description,
	}

	switch field.TypeAsString {
	case "Text", "Note", "ContentTypeId":
		c.Type = schema.TypeString
	case "Integer", "Counter":
		c.Type = schema.TypeInt
	case "Currency":
		c.Type = schema.TypeFloat
	case "Number":
		c.Type = schema.TypeFloat
	case "DateTime":
		c.Type = schema.TypeTimestamp
	case "Boolean", "Attachments":
		c.Type = schema.TypeBool
	case "Guid":
		c.Type = schema.TypeUUID
	case "Lookup", "User":
		c.Type = schema.TypeInt
	case "LookupMulti", "UserMulti":
		c.Type = schema.TypeIntArray
	case "Choice":
		c.Type = schema.TypeString
	case "MultiChoice":
		c.Type = schema.TypeStringArray
	case "Computed":
		c.Type = schema.TypeString
	default:
		logger.Warn().Str("type", field.TypeAsString).Int("kind", field.FieldTypeKind).Str("field_title", field.Title).Str("field_id", field.ID).Msg("unknown type, assuming JSON")
		c.Type = schema.TypeString
	}

	c.Name = normalizeName(field.InternalName)

	return c
}

func normalizeName(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.Trim(s, "_")
	return s
}
