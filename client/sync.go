package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
)

func (c *Client) Sync(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for _, table := range c.Tables {
		if metrics.TableClient[table.Name] == nil {
			metrics.TableClient[table.Name] = make(map[string]*source.TableClientMetrics)
			metrics.TableClient[table.Name][c.ID()] = &source.TableClientMetrics{}
		}
	}

	// Lists sync
	for tableName := range c.lists.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.lists.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}

	// MMD (Terms from TermSets) sync
	for tableName := range c.mmd.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.mmd.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}

	// User profiles sync
	for tableName := range c.profiles.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.profiles.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}

	// Search queries sync
	for tableName := range c.search.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.search.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}

	return nil
}
