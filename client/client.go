package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/rs/zerolog"
)

type Client struct {
	Tables schema.Tables

	lists *lists.Lists

	source specs.Source
	opts   source.Options
}

func (c *Client) ID() string {
	return c.source.Name
}

func NewClient(_ context.Context, logger zerolog.Logger, src specs.Source, opts source.Options) (schema.ClientMeta, error) {
	spec, err := getSpec(src)
	if err != nil {
		return nil, err
	}

	sp, err := auth.GetSP(spec.Auth)
	if err != nil {
		return nil, err
	}

	// sp.Conf(&api.RequestConfig{Context: ctx}) // for some reason gets context cancelled immediately

	client := &Client{
		lists:  lists.NewLists(sp, logger),
		source: src,
		opts:   opts,
	}

	client.Tables = make(schema.Tables, 0, len(spec.Lists))

	for listURI, listSpec := range spec.Lists {
		table, err := client.lists.GetDestTable(listURI, listSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from list: %w", err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("list", listURI).Str("columns", table.Columns.String()).Msg("columns for table")
			client.Tables = append(client.Tables, table)
		}
	}

	return client, nil
}
