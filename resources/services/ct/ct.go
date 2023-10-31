package ct

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/types"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type ContentTypesRollup struct {
	sp     *api.SP
	logger zerolog.Logger
}

func NewContentTypesRollup(sp *api.SP, logger zerolog.Logger) *ContentTypesRollup {
	return &ContentTypesRollup{
		sp:     sp,
		logger: logger,
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

	// ToDo: Rearchitect table construction logic
	for _, prop := range spec.Select {
		col := c.getDestCol(prop, tableName, ctInfo, spec)
		table.Columns = append(table.Columns, col)
	}

	table.Resolver = c.Resolver(ctInfo.ID, spec, table)

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

	valueResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
		value := util.GetRespValByProp(resource.Item.(map[string]interface{}), prop)
		if c.Type == arrow.BinaryTypes.String {
			if value != nil {
				value = fmt.Sprintf("%v", value)
			}
		}
		resource.Set(c.Name, value)
		return nil
	}

	// Props is not presented in list's fields
	if field == nil {
		return schema.Column{
			Name:        util.NormalizeEntityName(fieldAlias),
			Description: prop,
			Type:        c.typeFromPropName(prop),
			Resolver:    valueResolver,
		}
	}

	field.InternalName = fieldAlias
	col := c.columnFromField(field, tableName)
	col.Description = prop
	col.Resolver = valueResolver

	if prop == "UniqueId" {
		col.PrimaryKey = true
		col.Type = types.UUID
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
	case "ContentTypeId":
		col.Type = arrow.BinaryTypes.String
	case "Text", "Note":
		col.Type = arrow.BinaryTypes.LargeString
	case "Integer", "Counter":
		col.Type = arrow.PrimitiveTypes.Int32
	case "Currency":
		col.Type = arrow.PrimitiveTypes.Float32
	case "Number":
		col.Type = arrow.PrimitiveTypes.Float32
	case "DateTime":
		col.Type = arrow.FixedWidthTypes.Timestamp_us
	case "Boolean", "Attachments":
		col.Type = arrow.FixedWidthTypes.Boolean
	case "Guid":
		col.Type = types.UUID
	case "Lookup", "User":
		col.Type = arrow.PrimitiveTypes.Int32
	case "LookupMulti", "UserMulti":
		col.Type = arrow.ListOf(arrow.PrimitiveTypes.Int32)
	case "Choice":
		col.Type = arrow.BinaryTypes.String
	case "MultiChoice":
		col.Type = arrow.ListOf(arrow.BinaryTypes.String)
	case "Computed":
		col.Type = arrow.BinaryTypes.String
	default:
		logger.Warn().Str("type", field.TypeAsString).Int("kind", field.FieldTypeKind).Str("field_title", field.Title).Str("field_id", field.ID).Msg("unknown type, assuming JSON")
		col.Type = arrow.BinaryTypes.String
	}

	col.Name = util.NormalizeEntityName(field.InternalName)

	return col
}

func (*ContentTypesRollup) typeFromPropName(prop string) arrow.DataType {
	if strings.HasSuffix(prop, "/Id") && prop != "ParentList/Id" {
		return arrow.PrimitiveTypes.Int32
	}
	switch prop {
	case "ID", "Id", "AuthorId", "EditorId":
		return arrow.PrimitiveTypes.Int32
	case "ParentList/Id":
		return types.UUID
	case "ParentList/ParentWebUrl":
		return arrow.BinaryTypes.String
	case "Created", "Modified":
		return arrow.FixedWidthTypes.Timestamp_us
	default:
		return arrow.BinaryTypes.String
	}
}
