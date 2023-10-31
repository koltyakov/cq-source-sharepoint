package lists

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

type Lists struct {
	sp     *api.SP
	logger zerolog.Logger
}

func NewLists(sp *api.SP, logger zerolog.Logger) *Lists {
	return &Lists{
		sp:     sp,
		logger: logger,
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

	// ToDo: Rearchitect table construction logic
	for _, prop := range spec.Select {
		col := l.getDestCol(prop, tableName, spec, fieldsData)
		table.Columns = append(table.Columns, col)
	}

	table.Resolver = l.Resolver(listURI, spec, table)

	return table, nil
}

func (l *Lists) getDestCol(prop string, tableName string, spec Spec, fieldsData []api.FieldResp) schema.Column {
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
			Type:        typeFromPropName(prop),
			Resolver:    valueResolver,
		}
	}

	field.InternalName = fieldAlias
	col := l.columnFromField(field, tableName)
	col.PrimaryKey = prop == "ID" // ToDo: Decide on ID cunstruction logic: use ID/UniqueID/Path+ID
	col.Description = prop
	col.Resolver = valueResolver

	return col
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
		c.Type = arrow.BinaryTypes.String
	case "Integer", "Counter":
		c.Type = arrow.PrimitiveTypes.Int32
	case "Currency":
		c.Type = arrow.PrimitiveTypes.Float32
	case "Number":
		c.Type = arrow.PrimitiveTypes.Float32
	case "DateTime":
		c.Type = arrow.FixedWidthTypes.Timestamp_us
	case "Boolean", "Attachments":
		c.Type = arrow.FixedWidthTypes.Boolean
	case "Guid":
		c.Type = types.UUID
	case "Lookup", "User":
		c.Type = arrow.PrimitiveTypes.Int32
	case "LookupMulti", "UserMulti":
		c.Type = arrow.ListOf(arrow.PrimitiveTypes.Int32)
	case "Choice":
		c.Type = arrow.BinaryTypes.String
	case "MultiChoice":
		c.Type = arrow.ListOf(arrow.BinaryTypes.String)
	case "Computed":
		c.Type = arrow.BinaryTypes.String
	default:
		logger.Warn().Str("type", field.TypeAsString).Int("kind", field.FieldTypeKind).Str("field_title", field.Title).Str("field_id", field.ID).Msg("unknown type, assuming JSON")
		c.Type = arrow.BinaryTypes.String
	}

	c.Name = util.NormalizeEntityName(field.InternalName)

	return c
}

func typeFromPropName(prop string) arrow.DataType {
	if strings.HasSuffix(prop, "/Id") && prop != "ParentList/Id" {
		return arrow.PrimitiveTypes.Int32
	}
	return arrow.BinaryTypes.String
}
