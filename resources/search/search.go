package search

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
)

type Search struct {
	sp     *api.SP
	logger zerolog.Logger

	TablesMap map[string]Model // normalized table name to table metadata (map[CQ Table Name]Model)
}

type Model struct {
	Spec Spec
}

func NewSearch(sp *api.SP, logger zerolog.Logger) *Search {
	return &Search{
		sp:        sp,
		logger:    logger,
		TablesMap: map[string]Model{},
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
		fieldType := schema.TypeString
		for _, p := range ss {
			if p.Key == prop {
				switch p.ValueType {
				case "Edm.String":
					fieldType = schema.TypeString
				case "Edm.Int32":
				case "Edm.Int64":
					fieldType = schema.TypeInt
				case "Edm.Double":
					fieldType = schema.TypeFloat
				case "Edm.Boolean":
					fieldType = schema.TypeBool
				case "Edm.DateTime":
					fieldType = schema.TypeTimestamp
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
			{Name: "id", Type: schema.TypeInt, Description: "DocId", CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: util.NormalizeEntityNameSnake(getFieldAlias("Title", spec.fieldsMapping)), Type: schema.TypeString, Description: "Title"},
		}, columns...),
	}

	s.TablesMap[table.Name] = Model{
		Spec: spec,
	}

	return table, nil
}

func (s *Search) schemaBySpec(spec Spec) ([]*api.TypedKeyValue, error) {
	res, err := searchData(s.sp, spec, 0, 1)
	if err != nil {
		return nil, err
	}

	return res.Data().PrimaryQueryResult.RelevantResults.Table.Rows[0].Cells, nil
}

func getFieldAlias(field string, mapping map[string]string) string {
	if a, ok := mapping[field]; ok {
		return a
	}
	return field
}
