package profiles

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
)

// There is no way for actually getting all user profiles using client side APIs
// so we're mimicing this by getting users via Search API, so search should be up and running for this to work

type Profiles struct {
	sp     *api.SP
	logger zerolog.Logger

	TablesMap map[string]Model // normalized table name to table metadata (map[CQ Table Name]Model)
}

type Model struct {
	Spec      Spec
	FieldsMap map[string]string // cq column name to column metadata
}

func NewProfiles(sp *api.SP, logger zerolog.Logger) *Profiles {
	return &Profiles{
		sp:        sp,
		logger:    logger,
		TablesMap: map[string]Model{},
	}
}

var userProps = []string{"UniqueId", "Title", "WorkEmail", "JobTitle", "Department", "PictureURL", "AccountName", "Path", "LastModifiedTime"}

func (u *Profiles) GetDestTable(spec Spec) (*schema.Table, error) {
	tableName := "profile"
	if spec.Alias != "" {
		tableName = util.NormalizeEntityName(spec.Alias)
	}

	table := &schema.Table{
		Name:        "sharepoint_ups_" + tableName,
		Description: "User Profiles",
		Columns: []schema.Column{
			{Name: "id", Type: schema.TypeUUID, CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "title", Type: schema.TypeString},
			{Name: "email", Type: schema.TypeString},
			{Name: "job", Type: schema.TypeString},
			{Name: "department", Type: schema.TypeString},
			{Name: "picture", Type: schema.TypeString},
			{Name: "account", Type: schema.TypeString},
			{Name: "path", Type: schema.TypeString},
			{Name: "modified", Type: schema.TypeTimestamp},
		},
	}

	// ToDo: Remove this reverce mapping
	u.TablesMap[table.Name] = Model{
		Spec: spec,
		FieldsMap: map[string]string{
			"id":         "UniqueId",
			"title":      "Title",
			"email":      "WorkEmail",
			"job":        "JobTitle",
			"department": "Department",
			"picture":    "PictureURL",
			"account":    "AccountName",
			"path":       "Path",
			"modified":   "LastModifiedTime",
		},
	}

	return table, nil
}
