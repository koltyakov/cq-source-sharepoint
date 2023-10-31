package search

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
)

type Search struct {
	sp     *api.SP
	logger zerolog.Logger
}

func NewSearch(sp *api.SP, logger zerolog.Logger) *Search {
	return &Search{
		sp:     sp,
		logger: logger,
	}
}

func (s *Search) GetDestTable(searchName string, spec Spec) (*schema.Table, error) {
	tableName := util.NormalizeEntityName(searchName)

	ss, err := s.schemaBySpec(spec)
	if err != nil {
		return nil, err
	}

	columns := []schema.Column{}
	ignoreFields := []string{"DocId", "Title"}
	for _, prop := range spec.SelectProperties {
		var fieldType arrow.DataType = arrow.BinaryTypes.String
		for _, p := range ss {
			if p.Key == prop {
				switch p.ValueType {
				case "Edm.String":
					fieldType = arrow.BinaryTypes.String
				case "Edm.Int32":
					fieldType = arrow.PrimitiveTypes.Int32
				case "Edm.Int64":
					fieldType = arrow.PrimitiveTypes.Int64
				case "Edm.Double":
					fieldType = arrow.PrimitiveTypes.Float32
				case "Edm.Boolean":
					fieldType = arrow.FixedWidthTypes.Boolean
				case "Edm.DateTime":
					fieldType = arrow.FixedWidthTypes.Timestamp_us
				}
			}
		}

		if !util.Contains(ignoreFields, prop) {
			columns = append(columns, schema.Column{
				Name:        util.NormalizeEntityNameSnake(getFieldAlias(prop, spec.fieldsMapping)),
				Type:        fieldType,
				Description: prop,
			})
		}
	}

	table := &schema.Table{
		Name: "sharepoint_search_" + tableName,
		Columns: append([]schema.Column{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Description: "DocId", PrimaryKey: true},
			{Name: util.NormalizeEntityNameSnake(getFieldAlias("Title", spec.fieldsMapping)), Type: arrow.BinaryTypes.String, Description: "Title"},
		}, columns...),
	}

	for i, col := range table.Columns {
		prop := col.Description
		valueResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			value := getSearchCellValue(resource.Item.(*struct {
				Cells []*api.TypedKeyValue `json:"Cells"`
			}), prop)
			if c.Type == arrow.BinaryTypes.String {
				if value != nil {
					value = fmt.Sprintf("%v", value)
				}
			}
			return resource.Set(c.Name, value)
		}
		col.Resolver = valueResolver
		table.Columns[i] = col
	}

	table.Resolver = s.Resolver(spec, table)

	return table, nil
}

func (s *Search) schemaBySpec(spec Spec) ([]*api.TypedKeyValue, error) {
	res, err := searchData(s.sp, spec, 0, 1)
	if err != nil {
		return nil, err
	}

	rows := res.Data().PrimaryQueryResult.RelevantResults.Table.Rows

	if len(rows) == 0 {
		return []*api.TypedKeyValue{}, nil
	}

	return rows[0].Cells, nil
}

func getFieldAlias(field string, mapping map[string]string) string {
	if a, ok := mapping[field]; ok {
		return a
	}
	return field
}
