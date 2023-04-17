package ct

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type ContentTypesRollup struct {
	sp     *api.SP
	logger zerolog.Logger

	TablesMap map[string]Model // normalized table name to table metadata (map[CQ Table Name]Model)
}

type Model struct {
	ContentTypeID string
	Spec          Spec
}

func NewContentTypesRollup(sp *api.SP, logger zerolog.Logger) *ContentTypesRollup {
	return &ContentTypesRollup{
		sp:        sp,
		logger:    logger,
		TablesMap: map[string]Model{},
	}
}

func (c *ContentTypesRollup) GetDestTable(ctID string, spec Spec) (*schema.Table, error) {
	ctInfo, err := c.getContentTypeInfo(ctID)
	if err != nil {
		// ToDo: Decide which design is better to warn and go next or fail a sync
		// Will stay with a fast fail strateg for now so a user will know about an error immediately
		// Otherwise the only way to know about an error is to check `cloudquery.log` for
		// `2023-02-26T15:24:36Z ERR content type not found, skipping list={ctID} module=sharepoint-src`
		if util.IsNotFound(err) { // List not found, warn and skip
			c.logger.Error().Str("contentType", ctID).Msg("content type not found")
			return nil, fmt.Errorf("content type not found \"%s\": %w", ctID, err)
		}
		return nil, err
	}

	tableName := util.NormalizeEntityName(ctInfo.Name)
	if spec.Alias != "" {
		tableName = util.NormalizeEntityName(spec.Alias)
	}

	table := &schema.Table{
		Name:        "sharepoint_rollup_" + tableName,
		Description: ctInfo.Description,
	}

	model := &Model{
		ContentTypeID: ctInfo.ID,
		Spec:          spec,
	}

	// ToDo: Rearchitect table construction logic
	for _, prop := range spec.Select {
		col := c.getDestCol(prop, tableName, ctInfo, spec)
		table.Columns = append(table.Columns, col)
	}

	c.TablesMap[table.Name] = *model

	return table, nil
}

func (c *ContentTypesRollup) getDestCol(prop string, tableName string, ctInfo *contentTypeInfo, spec Spec) schema.Column {
	var field *api.FieldInfo
	for _, fieldData := range ctInfo.Fields {
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

	fieldAlias := prop
	if a, ok := spec.fieldsMapping[prop]; ok {
		fieldAlias = a
	}

	// Props is not presented in list's fields
	if field == nil {
		return schema.Column{
			Name:        util.NormalizeEntityName(fieldAlias),
			Description: prop,
			Type:        c.typeFromPropName(prop),
		}
	}

	field.InternalName = fieldAlias
	col := c.columnFromField(field, tableName)
	col.Description = prop

	if prop == "UniqueId" {
		col.CreationOptions.PrimaryKey = true
		col.Type = schema.TypeUUID
	}

	return col
}

type contentTypeInfo struct {
	ID          string           `json:"StringId"`
	Name        string           `json:"Name"`
	Description string           `json:"Description"`
	Fields      []*api.FieldInfo `json:"Fields"`
}

func (c *ContentTypesRollup) getContentTypeInfo(ctID string) (*contentTypeInfo, error) {
	resp, err := c.sp.Web().ContentTypes().
		Filter(fmt.Sprintf("Name eq '%s' or StringId eq '%s'", ctID, ctID)).
		Select("StringId,Name,Description").
		Expand("Fields").
		Top(5000).
		Get()

	if err != nil {
		return nil, err
	}

	var info []*contentTypeInfo
	if err := json.Unmarshal(resp.Normalized(), &info); err != nil {
		return nil, err
	}

	if len(info) == 0 {
		return nil, fmt.Errorf("content type not found: %s", ctID)
	}

	return info[0], nil
}

func (c *ContentTypesRollup) columnFromField(field *api.FieldInfo, tableName string) schema.Column {
	logger := c.logger.With().Str("table", tableName).Logger()

	col := schema.Column{
		Description: field.Description,
	}

	switch field.TypeAsString {
	case "Text", "Note", "ContentTypeId":
		col.Type = schema.TypeString
	case "Integer", "Counter":
		col.Type = schema.TypeInt
	case "Currency":
		col.Type = schema.TypeFloat
	case "Number":
		col.Type = schema.TypeFloat
	case "DateTime":
		col.Type = schema.TypeTimestamp
	case "Boolean", "Attachments":
		col.Type = schema.TypeBool
	case "Guid":
		col.Type = schema.TypeUUID
	case "Lookup", "User":
		col.Type = schema.TypeInt
	case "LookupMulti", "UserMulti":
		col.Type = schema.TypeIntArray
	case "Choice":
		col.Type = schema.TypeString
	case "MultiChoice":
		col.Type = schema.TypeStringArray
	case "Computed":
		col.Type = schema.TypeString
	default:
		logger.Warn().Str("type", field.TypeAsString).Int("kind", field.FieldTypeKind).Str("field_title", field.Title).Str("field_id", field.ID).Msg("unknown type, assuming JSON")
		col.Type = schema.TypeString
	}

	col.Name = util.NormalizeEntityName(field.InternalName)

	return col
}

func (*ContentTypesRollup) typeFromPropName(prop string) schema.ValueType {
	if strings.HasSuffix(prop, "/Id") && prop != "ParentList/Id" {
		return schema.TypeInt
	}
	switch prop {
	case "ID", "Id", "AuthorId", "EditorId":
		return schema.TypeInt
	case "ParentList/Id":
		return schema.TypeUUID
	case "ParentList/ParentWebUrl":
		return schema.TypeString
	case "Created", "Modified":
		return schema.TypeTimestamp
	default:
		return schema.TypeString
	}
}
