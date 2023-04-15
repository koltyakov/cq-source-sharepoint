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

	if err := c.syncLists(ctx, metrics, res); err != nil {
		return err
	}

	if err := c.syncMMD(ctx, metrics, res); err != nil {
		return err
	}

	if err := c.syncProfiles(ctx, metrics, res); err != nil {
		return err
	}

	if err := c.syncSearch(ctx, metrics, res); err != nil {
		return err
	}

	if err := c.syncContentTypes(ctx, metrics, res); err != nil {
		return err
	}

	return nil
}

// Lists sync
func (c *Client) syncLists(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for tableName := range c.lists.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.lists.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// MMD (Terms from TermSets) sync
func (c *Client) syncMMD(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for tableName := range c.mmd.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.mmd.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// User profiles sync
func (c *Client) syncProfiles(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for tableName := range c.profiles.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.profiles.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// Search queries sync
func (c *Client) syncSearch(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for tableName := range c.search.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.search.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

// Content types rollup sync
func (c *Client) syncContentTypes(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for tableName := range c.contentTypes.TablesMap {
		table := c.Tables.Get(tableName)
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.contentTypes.Sync(ctx, m, res, table); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}
