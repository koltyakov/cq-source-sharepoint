package lists

import (
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type Lists struct {
	sp     *api.SP
	logger zerolog.Logger

	TablesMap map[string]Model // normalized table name to table metadata (map[CQ Table Name]Model)
}

type Model struct {
	URI  string
	Spec Spec
}

func NewLists(sp *api.SP, logger zerolog.Logger) *Lists {
	return &Lists{
		sp:        sp,
		logger:    logger,
		TablesMap: map[string]Model{},
	}
}

func (l *Lists) GetDestTable(listURI string, spec Spec) (*schema.Table, error) {
	listInfo, err := l.getListInfo(listURI)
	if err != nil {
		// ToDo: Decide which design is better to warn and go next or fail a sync
		// Will stay with a fast fail strateg for now so a user will know about an error immediately
		// Otherwise the only way to know about an error is to check `cloudquery.log` for
		// `2023-02-26T15:24:36Z ERR list not found, skipping list={ListURI} module=sharepoint-src`
		if util.IsNotFound(err) { // List not found, warn and skip
			l.logger.Error().Str("list", listURI).Msg("list not found")
			return nil, fmt.Errorf("list not found \"%s\": %w", listURI, err)
		}
		return nil, err
	}

	siteURL := util.GetRelativeURL(l.sp.ToURL())
	lURI := util.RemoveRelativeURLPrefix(listInfo.RootFolder.ServerRelativeURL, siteURL)

	tableName := util.NormalizeEntityName(lURI)
	if spec.Alias != "" {
		tableName = util.NormalizeEntityName(spec.Alias)
	}

	table := &schema.Table{
		Name:        "sharepoint_" + tableName,
		Description: listInfo.Description,
	}

	fields, err := l.sp.Web().GetList(listURI).Fields().Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}

	fieldsData := fields.Data()
	model := &Model{
		URI:  listURI,
		Spec: spec,
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
				Name:        util.NormalizeEntityName(prop),
				Description: prop,
				Type:        schema.TypeString,
			}

			table.Columns = append(table.Columns, c)
			continue
		}

		col := l.columnFromField(field, table.Name)
		col.CreationOptions.PrimaryKey = prop == "ID" // ToDo: Decide on ID cunstruction logic: use ID/UniqueID/Path+ID
		col.Description = prop

		table.Columns = append(table.Columns, col)
	}

	l.TablesMap[table.Name] = *model

	return table, nil
}

type listInfo struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	RootFolder  struct {
		ServerRelativeURL string `json:"ServerRelativeUrl"`
	} `json:"RootFolder"`
}

func (l *Lists) getListInfo(listURI string) (*listInfo, error) {
	list := l.sp.Web().GetList(listURI)

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

func (l *Lists) columnFromField(field *api.FieldInfo, tableName string) schema.Column {
	logger := l.logger.With().Str("table", tableName).Logger()

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

	c.Name = util.NormalizeEntityName(field.InternalName)

	return c
}
