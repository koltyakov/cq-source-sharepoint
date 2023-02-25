package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/thoas/go-funk"
)

func (c *Client) Sync(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for _, table := range c.Tables {
		if metrics.TableClient[table.Name] == nil {
			metrics.TableClient[table.Name] = make(map[string]*source.TableClientMetrics)
			metrics.TableClient[table.Name][c.ID()] = &source.TableClientMetrics{}
		}
	}

	for _, table := range c.Tables {
		meta := c.tablesMap[table.Name]
		m := metrics.TableClient[table.Name][c.ID()]
		if err := c.syncTable(ctx, m, res, table, meta); err != nil {
			return fmt.Errorf("syncing table %s: %w", table.Name, err)
		}
	}
	return nil
}

func (c *Client) syncTable(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table, meta ListModel) error {
	logger := c.Logger.With().Str("table", table.Name).Logger()

	logger.Debug().Strs("cols", meta.ListSpec.Select).Msg("selecting columns from list")

	list := c.SP.Web().GetList(meta.ListURI)
	items, err := list.Items().
		Select(strings.Join(meta.ListSpec.Select, ",")).
		Expand(strings.Join(meta.ListSpec.Expand, ",")).
		Top(2000).GetPaged()

	for {
		if err != nil {
			metrics.Errors++
			return fmt.Errorf("failed to get items: %w", err)
		}

		var itemList []map[string]any
		if err := json.Unmarshal(items.Items.Normalized(), &itemList); err != nil {
			metrics.Errors++
			return err
		}

		for _, itemMap := range itemList {
			ks := funk.Keys(itemMap).([]string)
			sort.Strings(ks)
			logger.Debug().Strs("keys", ks).Msg("item keys")

			colVals := make([]any, len(table.Columns))

			for i, col := range table.Columns {
				prop := meta.FieldsMap[col.Name]
				colVals[i] = getRespValByProp(itemMap, prop)
			}

			resource, err := resourceFromValues(table, colVals)
			if err != nil {
				metrics.Errors++
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case res <- resource:
				metrics.Resources++
			}
		}

		if !items.HasNextPage() {
			break
		}
		items, err = items.GetNextPage()
	}

	return nil
}

func resourceFromValues(table *schema.Table, values []any) (*schema.Resource, error) {
	resource := schema.NewResourceData(table, nil, values)
	for i, col := range table.Columns {
		if err := resource.Set(col.Name, values[i]); err != nil {
			return nil, err
		}
	}
	return resource, nil
}
