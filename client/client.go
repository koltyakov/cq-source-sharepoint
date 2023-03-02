package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/search"
	"github.com/rs/zerolog"
)

type Client struct {
	Tables schema.Tables

	lists    *lists.Lists
	mmd      *mmd.MMD
	profiles *profiles.Profiles
	search   *search.Search

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
		lists:    lists.NewLists(sp, logger),
		mmd:      mmd.NewMMD(sp, logger),
		profiles: profiles.NewProfiles(sp, logger),
		search:   search.NewSearch(sp, logger),

		source: src,
		opts:   opts,
	}

	tableCnt := len(spec.Lists) + len(spec.MMD)
	if spec.Profiles.Enabled {
		tableCnt++
	}
	client.Tables = make(schema.Tables, 0, tableCnt)

	// Managed metadata tables prepare
	for termSetID, mmdSpec := range spec.MMD {
		table, err := client.mmd.GetDestTable(termSetID, mmdSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from term set \"%s\": %w", termSetID, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("termset", termSetID).Str("columns", table.Columns.String()).Msg("columns for table")
			client.Tables = append(client.Tables, table)
		}
	}

	// Lists tables prepare
	for listURI, listSpec := range spec.Lists {
		table, err := client.lists.GetDestTable(listURI, listSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from list \"%s\": %w", listURI, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("list", listURI).Str("columns", table.Columns.String()).Msg("columns for table")
			client.Tables = append(client.Tables, table)
		}
	}

	// User profiles tables prepare
	if spec.Profiles.Enabled {
		table, err := client.profiles.GetDestTable(spec.Profiles)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from user profiles: %w", err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("columns", table.Columns.String()).Msg("columns for table")
			client.Tables = append(client.Tables, table)
		}
	}

	// Search tables prepare
	for searchName, searchSpec := range spec.Search {
		table, err := client.search.GetDestTable(searchName, searchSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from search query \"%s\": %w", searchName, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("search", searchName).Str("columns", table.Columns.String()).Msg("columns for table")
			client.Tables = append(client.Tables, table)
		}
	}

	return client, nil
}
