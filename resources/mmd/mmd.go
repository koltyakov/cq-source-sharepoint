package mmd

import (
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
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
	ID        string
	Spec      Spec
	FieldsMap map[string]string // cq column name to column metadata
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
			{Name: "id", Type: schema.TypeUUID, CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "name", Type: schema.TypeString},
			{Name: "description", Type: schema.TypeString},
			{Name: "tagging", Type: schema.TypeBool},
			{Name: "deprecated", Type: schema.TypeBool},
			{Name: "pinned", Type: schema.TypeBool},
			{Name: "reused", Type: schema.TypeBool},
			{Name: "root", Type: schema.TypeBool},
			{Name: "source", Type: schema.TypeBool},
			{Name: "path", Type: schema.TypeStringArray},
			{Name: "children", Type: schema.TypeInt},
			{Name: "merged", Type: schema.TypeUUIDArray},
			{Name: "shared_props", Type: schema.TypeJSON},
			{Name: "local_props", Type: schema.TypeJSON},
			{Name: "custom_sort", Type: schema.TypeUUIDArray},
			{Name: "owner", Type: schema.TypeString},
			{Name: "created", Type: schema.TypeTimestamp},
			{Name: "modified", Type: schema.TypeTimestamp},
		},
	}

	// ToDo: Remove this reverce mapping
	m.TablesMap[table.Name] = Model{
		ID:   terSetID,
		Spec: spec,
		FieldsMap: map[string]string{
			"id":           "Id",
			"name":         "Name",
			"description":  "Description",
			"tagging":      "IsAvailableForTagging",
			"deprecated":   "IsDeprecated",
			"pinned":       "IsPinned",
			"reused":       "IsReused",
			"root":         "IsRoot",
			"source":       "IsSourceTerm",
			"path":         "PathOfTerm",
			"children":     "TermsCount",
			"merged":       "MergedTermIds",
			"shared_props": "CustomProperties",
			"local_props":  "LocalCustomProperties",
			"custom_sort":  "CustomSortOrder",
			"owner":        "Owner",
			"created":      "CreatedDate",
			"modified":     "LastModifiedDate",
		},
	}

	return table, nil
}
