package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/caser"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/addin"
	"github.com/rs/zerolog"
)

type Client struct {
	Logger     zerolog.Logger
	Tables     schema.Tables
	SP         *api.SP
	spec       specs.Source
	pluginSpec Spec
	opts       source.Options
	csr        *caser.Caser

	tablesMap map[string]tableMeta // normalized table name to table metadata
}

type tableMeta struct {
	Title     string
	ColumnMap map[string]columnMeta // cq column name to column metadata
}

type columnMeta struct {
	SharepointName string
	SharepointType string
}

func (c *Client) ID() string {
	return c.spec.Name
}

func New(_ context.Context, logger zerolog.Logger, s specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var pluginSpec Spec

	if err := s.UnmarshalSpec(&pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}
	pluginSpec.SetDefaults()
	if err := pluginSpec.Validate(); err != nil {
		return nil, err
	}

	auth := &strategy.AuthCnfg{
		SiteURL:      pluginSpec.SiteURL,
		ClientID:     pluginSpec.ClientID,
		ClientSecret: pluginSpec.ClientSecret,
	}
	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)

	cl := &Client{
		Logger:     logger,
		SP:         sp,
		spec:       s,
		pluginSpec: pluginSpec,
		opts:       opts,
		csr:        caser.New(),
	}

	// if len(pluginSpec.Lists) == 0 {
	// 	var err error
	// 	pluginSpec.Lists, err = cl.getAllLists()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	cl.Tables = make(schema.Tables, 0, len(pluginSpec.Lists))
	cl.tablesMap = make(map[string]tableMeta, len(pluginSpec.Lists))
	for _, title := range pluginSpec.Lists {
		table, meta, err := cl.tableFromList(title)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from list: %w", err)
		}
		if table != nil {
			// b, _ := json.Marshal(meta.ColumnMap)
			// fmt.Println(string(b))
			cl.Logger.Debug().Str("table", table.Name).Str("list", title).Str("columns", table.Columns.String()).Msg("columns for table")

			cl.Tables = append(cl.Tables, table)
			cl.tablesMap[table.Name] = *meta
		}
	}

	return cl, nil
}
