package plugin

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/ct"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/search"
	"github.com/rs/zerolog"
)

type Client struct {
	lists        *lists.Lists
	mmd          *mmd.MMD
	profiles     *profiles.Profiles
	search       *search.Search
	contentTypes *ct.ContentTypesRollup
}

func NewClient(ctx context.Context, logger zerolog.Logger, cnfg []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	spec, err := getSpec(cnfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	sp, err := auth.GetSP(spec.Auth)
	if err != nil {
		return nil, err
	}

	if _, err := sp.Web().Select("Title").Get(); err != nil {
		return nil, fmt.Errorf("failed to connect to SharePoint: %w", err)
	}

	client := &Client{
		lists:        lists.NewLists(sp, logger),
		mmd:          mmd.NewMMD(sp, logger),
		profiles:     profiles.NewProfiles(sp, logger),
		search:       search.NewSearch(sp, logger),
		contentTypes: ct.NewContentTypesRollup(sp, logger),
	}

	tables, err := client.getTables(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tables: %w", err)
	}

	if opts.NoConnection {
		return &Plugin{
			logger: logger,
			tables: tables,
		}, nil
	}

	return &Plugin{
		logger:    logger,
		spec:      *spec,
		tables:    tables,
		scheduler: scheduler.NewScheduler(scheduler.WithLogger(logger)),
		client:    client,
	}, nil
}

func (c *Client) getTables(config *Spec) (schema.Tables, error) {
	tables := schema.Tables{}

	// Tables from lists config
	for listURI, listSpec := range config.Lists {
		table, err := c.lists.GetDestTable(listURI, listSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get list '%s': %w", listURI, err)
		}
		tables = append(tables, table)
	}

	// Tables from mmd config
	for terSetID, mmdSpec := range config.MMD {
		table, err := c.mmd.GetDestTable(terSetID, mmdSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get term set '%s': %w", terSetID, err)
		}
		tables = append(tables, table)
	}

	// Tables from profiles config
	if config.Profiles.Enabled {
		table, err := c.profiles.GetDestTable(config.Profiles)
		if err != nil {
			return nil, fmt.Errorf("failed to get profiles: %w", err)
		}
		tables = append(tables, table)
	}

	// Tables from search config
	for searchName, searchSpec := range config.Search {
		table, err := c.search.GetDestTable(searchName, searchSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get search '%s': %w", searchName, err)
		}
		tables = append(tables, table)
	}

	// Tables from content types config
	for ctName, ctSpec := range config.ContentTypes {
		table, err := c.contentTypes.GetDestTable(ctName, ctSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to get content type '%s': %w", ctName, err)
		}
		tables = append(tables, table)
	}

	if err := transformers.TransformTables(tables); err != nil {
		return nil, err
	}

	for _, table := range tables {
		schema.AddCqIDs(table)
	}

	return tables, nil
}
