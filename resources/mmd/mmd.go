package mmd

import (
	"strings"

	"github.com/cloudquery/plugin-sdk/v2/schema"
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
			{Name: "id", Type: schema.TypeUUID, Description: "Id", CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "name", Type: schema.TypeString, Description: "Name"},
			{Name: "description", Type: schema.TypeString, Description: "Description"},
			{Name: "tagging", Type: schema.TypeBool, Description: "IsAvailableForTagging"},
			{Name: "deprecated", Type: schema.TypeBool, Description: "IsDeprecated"},
			{Name: "pinned", Type: schema.TypeBool, Description: "IsPinned"},
			{Name: "reused", Type: schema.TypeBool, Description: "IsReused"},
			{Name: "root", Type: schema.TypeBool, Description: "IsRoot"},
			{Name: "source", Type: schema.TypeBool, Description: "IsSourceTerm"},
			{Name: "path", Type: schema.TypeStringArray, Description: "Path"},
			{Name: "children", Type: schema.TypeInt, Description: "ChildrenCount"},
			{Name: "merged", Type: schema.TypeUUIDArray, Description: "MergedTermIds"},
			{Name: "shared_props", Type: schema.TypeJSON, Description: "CustomProperties"},
			{Name: "local_props", Type: schema.TypeJSON, Description: "LocalCustomProperties"},
			{Name: "custom_sort", Type: schema.TypeUUIDArray, Description: "CustomSortOrder"},
			{Name: "owner", Type: schema.TypeString, Description: "Owner"},
			{Name: "created", Type: schema.TypeTimestamp, Description: "CreatedDate"},
			{Name: "modified", Type: schema.TypeTimestamp, Description: "LastModifiedDate"},
		},
	}

	m.TablesMap[table.Name] = Model{
		ID:   terSetID,
		Spec: spec,
	}

	return table, nil
}
