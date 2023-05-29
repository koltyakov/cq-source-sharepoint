package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/source"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/ct"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/search"
	"github.com/rs/zerolog"
)

type Client struct {
	Tables schema.Tables

	lists        *lists.Lists
	mmd          *mmd.MMD
	profiles     *profiles.Profiles
	search       *search.Search
	contentTypes *ct.ContentTypesRollup

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
		lists:        lists.NewLists(sp, logger),
		mmd:          mmd.NewMMD(sp, logger),
		profiles:     profiles.NewProfiles(sp, logger),
		search:       search.NewSearch(sp, logger),
		contentTypes: ct.NewContentTypesRollup(sp, logger),

		source: src,
		opts:   opts,
	}

	tableCnt := len(spec.Lists) + len(spec.MMD) + len(spec.Search) + len(spec.ContentTypes)
	if spec.Profiles.Enabled {
		tableCnt++
	}
	client.Tables = make(schema.Tables, 0, tableCnt)

	// Managed metadata tables prepare
	mmdTables, err := prepareMMDTables(client, spec, logger)
	if err != nil {
		return nil, err
	}
	client.Tables = append(client.Tables, mmdTables...)

	// Lists tables prepare
	listTables, err := prepareListTables(client, spec, logger)
	if err != nil {
		return nil, err
	}
	client.Tables = append(client.Tables, listTables...)

	// User profiles tables prepare
	profileTables, err := prepareProfileTables(client, spec, logger)
	if err != nil {
		return nil, err
	}
	client.Tables = append(client.Tables, profileTables...)

	// Search tables prepare
	searchTables, err := prepareSearchTables(client, spec, logger)
	if err != nil {
		return nil, err
	}
	client.Tables = append(client.Tables, searchTables...)

	// Content types rollup tables prepare
	ctTables, err := prepareContentTypeTables(client, spec, logger)
	if err != nil {
		return nil, err
	}
	client.Tables = append(client.Tables, ctTables...)

	return client, nil
}

func prepareMMDTables(client *Client, spec *Spec, logger zerolog.Logger) ([]*schema.Table, error) {
	tables := make([]*schema.Table, 0, len(spec.MMD))
	for termSetID, mmdSpec := range spec.MMD {
		table, err := client.mmd.GetDestTable(termSetID, mmdSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from term set \"%s\": %w", termSetID, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("termset", termSetID).Str("columns", table.Columns.String()).Msg("columns for table")
			tables = append(tables, table)
		}
	}
	return tables, nil
}

func prepareListTables(client *Client, spec *Spec, logger zerolog.Logger) ([]*schema.Table, error) {
	tables := make([]*schema.Table, 0, len(spec.Lists))
	for listURI, listSpec := range spec.Lists {
		table, err := client.lists.GetDestTable(listURI, listSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from list \"%s\": %w", listURI, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("list", listURI).Str("columns", table.Columns.String()).Msg("columns for table")
			tables = append(tables, table)
		}
	}
	return tables, nil
}

func prepareProfileTables(client *Client, spec *Spec, logger zerolog.Logger) ([]*schema.Table, error) {
	if spec.Profiles.Enabled {
		table, err := client.profiles.GetDestTable(spec.Profiles)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from user profiles: %w", err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("columns", table.Columns.String()).Msg("columns for table")
			return []*schema.Table{table}, nil
		}
	}
	return []*schema.Table{}, nil
}

func prepareSearchTables(client *Client, spec *Spec, logger zerolog.Logger) ([]*schema.Table, error) {
	tables := make([]*schema.Table, 0, len(spec.Search))
	for searchName, searchSpec := range spec.Search {
		table, err := client.search.GetDestTable(searchName, searchSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from search query \"%s\": %w", searchName, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("search", searchName).Str("columns", table.Columns.String()).Msg("columns for table")
			tables = append(tables, table)
		}
	}
	return tables, nil
}

func prepareContentTypeTables(client *Client, spec *Spec, logger zerolog.Logger) ([]*schema.Table, error) {
	tables := make([]*schema.Table, 0, len(spec.ContentTypes))
	for ctName, ctSpec := range spec.ContentTypes {
		table, err := client.contentTypes.GetDestTable(ctName, ctSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get table from content type \"%s\": %w", ctName, err)
		}
		if table != nil {
			logger.Debug().Str("table", table.Name).Str("contenttype", ctName).Str("columns", table.Columns.String()).Msg("columns for table")
			tables = append(tables, table)
		}
	}
	return tables, nil
}
