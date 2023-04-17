package profiles

import (
	"github.com/cloudquery/plugin-sdk/v2/schema"
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
	Spec Spec
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
			{Name: "id", Type: schema.TypeUUID, Description: "UniqueId", CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true}},
			{Name: "title", Type: schema.TypeString, Description: "Title"},
			{Name: "email", Type: schema.TypeString, Description: "WorkEmail"},
			{Name: "job", Type: schema.TypeString, Description: "JobTitle"},
			{Name: "department", Type: schema.TypeString, Description: "Department"},
			{Name: "picture", Type: schema.TypeString, Description: "PictureURL"},
			{Name: "account", Type: schema.TypeString, Description: "AccountName"},
			{Name: "path", Type: schema.TypeString, Description: "Path"},
			{Name: "modified", Type: schema.TypeTimestamp, Description: "LastModifiedTime"},
		},
	}

	u.TablesMap[table.Name] = Model{
		Spec: spec,
	}

	return table, nil
}
