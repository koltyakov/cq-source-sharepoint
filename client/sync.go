package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	// for _, table := range c.Tables {
	// 	if metrics.TableClient[table.Name] == nil {
	// 		metrics.TableClient[table.Name] = make(map[string]*source.TableClientMetrics)
	// 		metrics.TableClient[table.Name][c.ID()] = &source.TableClientMetrics{}
	// 	}
	// }

	if err := c.syncLists(ctx, options, res); err != nil {
		return err
	}

	if err := c.syncMMD(ctx, options, res); err != nil {
		return err
	}

	if err := c.syncProfiles(ctx, options, res); err != nil {
		return err
	}

	if err := c.syncSearch(ctx, options, res); err != nil {
		return err
	}

	return c.syncContentTypes(ctx, options, res)
}

func (*Client) Close(_ context.Context) error {
	// ToDo: Add your client cleanup here
	return nil
}

// Lists sync
func (c *Client) syncLists(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	for tableName := range c.lists.TablesMap {
		table := c.Tables.Get(tableName)
		// m := metrics.TableClient[table.Name][c.ID()]
		if err := c.lists.Sync(ctx, options, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// MMD (Terms from TermSets) sync
func (c *Client) syncMMD(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	for tableName := range c.mmd.TablesMap {
		table := c.Tables.Get(tableName)
		// m := metrics.TableClient[table.Name][c.ID()]
		if err := c.mmd.Sync(ctx, options, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// User profiles sync
func (c *Client) syncProfiles(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	for tableName := range c.profiles.TablesMap {
		table := c.Tables.Get(tableName)
		// m := metrics.TableClient[table.Name][c.ID()]
		if err := c.profiles.Sync(ctx, options, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// Search queries sync
func (c *Client) syncSearch(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	for tableName := range c.search.TablesMap {
		table := c.Tables.Get(tableName)
		// m := metrics.TableClient[table.Name][c.ID()]
		if err := c.search.Sync(ctx, options, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// Content types rollup sync
func (c *Client) syncContentTypes(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	for tableName := range c.contentTypes.TablesMap {
		table := c.Tables.Get(tableName)
		// m := metrics.TableClient[table.Name][c.ID()]
		if err := c.contentTypes.Sync(ctx, options, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}
