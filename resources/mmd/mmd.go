package mmd

import (
	"strings"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
)

type MMD struct {
	sp     *api.SP
	logger zerolog.Logger

	TablesMap map[string]Model // normalized table name to table metadata (map[CQ Table Name]Model)
}

type Model struct {
	ID   string
	Spec Spec
}

func NewMMD(sp *api.SP, logger zerolog.Logger) *MMD {
	return &MMD{
		sp:        sp,
		logger:    logger,
		TablesMap: map[string]Model{},
	}
}

func (m *MMD) GetDestTable(terSetID string, spec Spec) (*schema.Table, error) {
	tableName := util.NormalizeEntityName(strings.ReplaceAll(terSetID, "-", "")) // ToDo: ${TertGoup}_${TermSetName}
	if spec.Alias != "" {
		tableName = util.NormalizeEntityName(spec.Alias)
	}

	table := &schema.Table{
		Name:        "sharepoint_mmd_" + tableName,
		Description: "", // TermSetName
		Columns: []schema.Column{
			{Name: "id", Type: types.UUID, Description: "Id", PrimaryKey: true},
			{Name: "name", Type: arrow.BinaryTypes.String, Description: "Name"},
			{Name: "description", Type: arrow.BinaryTypes.String, Description: "Description"},
			{Name: "tagging", Type: arrow.FixedWidthTypes.Boolean, Description: "IsAvailableForTagging"},
			{Name: "deprecated", Type: arrow.FixedWidthTypes.Boolean, Description: "IsDeprecated"},
			{Name: "pinned", Type: arrow.FixedWidthTypes.Boolean, Description: "IsPinned"},
			{Name: "reused", Type: arrow.FixedWidthTypes.Boolean, Description: "IsReused"},
			{Name: "root", Type: arrow.FixedWidthTypes.Boolean, Description: "IsRoot"},
			{Name: "source", Type: arrow.FixedWidthTypes.Boolean, Description: "IsSourceTerm"},
			{Name: "path", Type: arrow.ListOf(arrow.BinaryTypes.String), Description: "Path"},
			{Name: "children", Type: arrow.PrimitiveTypes.Int32, Description: "ChildrenCount"},
			{Name: "merged", Type: arrow.ListOf(types.UUID), Description: "MergedTermIds"},
			{Name: "shared_props", Type: types.ExtensionTypes.JSON, Description: "CustomProperties"},
			{Name: "local_props", Type: types.ExtensionTypes.JSON, Description: "LocalCustomProperties"},
			{Name: "custom_sort", Type: arrow.ListOf(types.UUID), Description: "CustomSortOrder"},
			{Name: "owner", Type: arrow.BinaryTypes.String, Description: "Owner"},
			{Name: "created", Type: arrow.FixedWidthTypes.Timestamp_us, Description: "CreatedDate"},
			{Name: "modified", Type: arrow.FixedWidthTypes.Timestamp_us, Description: "LastModifiedDate"},
		},
	}

	m.TablesMap[table.Name] = Model{
		ID:   terSetID,
		Spec: spec,
	}

	return table, nil
}
