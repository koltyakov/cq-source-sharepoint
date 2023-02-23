package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/gosip/api"
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

func (c *Client) syncTable(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table, meta tableMeta) error {
	logger := c.Logger.With().Str("table", table.Name).Logger()

	colsToSelect := make([]string, 0, len(meta.ColumnMap))
	for _, v := range meta.ColumnMap {
		colsToSelect = append(colsToSelect, v.SharepointName)
	}
	logger.Debug().Strs("cols", colsToSelect).Msg("selecting columns from list")

	list := c.SP.Web().GetList("Lists/" + meta.Title)
	items, err := list.Items().Select(strings.Join(colsToSelect, ", ")).GetPaged()

	for {
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
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
			var notFoundCols []string

			for i, col := range table.Columns {
				colMeta := meta.ColumnMap[col.Name]
				val, ok := itemMap[colMeta.SharepointName]
				if !ok {
					notFoundCols = append(notFoundCols, colMeta.SharepointName)
					colVals[i] = nil
					continue
				}
				colVals[i] = convertSharepointType(colMeta, val)
				delete(itemMap, colMeta.SharepointName)
			}

			if len(notFoundCols) > 0 {
				sort.Strings(notFoundCols)
				logger.Warn().Strs("missing_columns", notFoundCols).Msg("missing columns in result")
			}
			if len(itemMap) > 0 {
				// Remove any extra fields that are already ignored but still found in the response
				for k := range itemMap {
					if !c.pluginSpec.ShouldSelectField(meta.Title, api.FieldInfo{InternalName: k}) {
						delete(itemMap, k)
					}
				}
				delete(itemMap, "Id") // remove "Id", we should already have "ID"
			}
			if len(itemMap) > 0 {
				ks := funk.Keys(itemMap).([]string)
				sort.Strings(ks)
				logger.Warn().Strs("extra_columns", ks).Msg("extra columns found in result")
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

func convertSharepointType(colMeta columnMeta, val any) any {
	switch colMeta.SharepointType {
	case "Currency":
		return fmt.Sprintf("%f", val)
	default:
		return val
	}
}
