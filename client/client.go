package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	"github.com/koltyakov/gosip/auth"
	"github.com/rs/zerolog"
)

type Client struct {
	Logger zerolog.Logger
	Tables schema.Tables
	SP     *api.SP

	src  specs.Source
	spec Spec
	opts source.Options

	tablesMap map[string]ListModel // normalized table name to table metadata
}

type ListModel struct {
	ListURI   string
	ListSpec  ListSpec
	FieldsMap map[string]string // cq column name to column metadata
}

func (c *Client) ID() string {
	return c.src.Name
}

func New(_ context.Context, logger zerolog.Logger, src specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var spec Spec

	if err := src.UnmarshalSpec(&spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}

	spec.SetDefaults()

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	jsonCreds, _ := json.Marshal(spec.Auth.Creds)
	authCnfg, err := auth.NewAuthByStrategy(spec.Auth.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth config: %w", err)
	}
	if err := authCnfg.ParseConfig(jsonCreds); err != nil {
		return nil, fmt.Errorf("failed to parse auth config: %w", err)
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	cl := &Client{
		Logger: logger,
		SP:     sp,

		src:  src,
		spec: spec,
		opts: opts,
	}

	cl.Tables = make(schema.Tables, 0, len(spec.Lists))
	cl.tablesMap = make(map[string]ListModel, len(spec.Lists))

	for listURI, listSpec := range spec.Lists {
		table, meta, err := cl.tableFromList(listURI, listSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from list: %w", err)
		}
		if table != nil {
			cl.Logger.Debug().Str("table", table.Name).Str("list", listURI).Str("columns", table.Columns.String()).Msg("columns for table")

			cl.Tables = append(cl.Tables, table)
			cl.tablesMap[table.Name] = *meta
		}
	}

	return cl, nil
}
