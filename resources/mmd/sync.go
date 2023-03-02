package mmd

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/thoas/go-funk"
)

func (m *MMD) Sync(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table) error {
	opts := m.TablesMap[table.Name]
	logger := m.logger.With().Str("table", table.Name).Logger()

	taxonomy := m.sp.Taxonomy()
	terms, err := taxonomy.Stores().Default().Sets().GetByID(opts.ID).GetAllTerms()
	if err != nil {
		metrics.Errors++
		return fmt.Errorf("failed to get items: %w", err)
	}

	for _, itemMap := range terms {
		ks := funk.Keys(itemMap).([]string)
		sort.Strings(ks)
		logger.Debug().Strs("keys", ks).Msg("item keys")

		colVals := make([]any, len(table.Columns))

		for i, col := range table.Columns {
			prop := col.Description
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

func getRespValByProp(resp map[string]any, prop string) any {
	val := util.GetRespValByProp(resp, prop)
	if prop == "Id" {
		return strings.ReplaceAll(strings.ReplaceAll(val.(string), "/Guid(", ""), ")/", "")
	}
	if prop == "CreatedDate" || prop == "LastModifiedDate" {
		dateStr := strings.ReplaceAll(strings.ReplaceAll(val.(string), "/Date(", ""), ")/", "")
		dateInt, _ := strconv.ParseInt(dateStr, 10, 64)
		return time.UnixMilli(dateInt)
	}
	if prop == "PathOfTerm" {
		return strings.Split(val.(string), ";")
	}
	if prop == "MergedTermIds" {
		mergedTerms := val.([]any)
		for i, term := range mergedTerms {
			mergedTerms[i] = strings.ReplaceAll(strings.ReplaceAll(term.(string), "/Guid(", ""), ")/", "")
		}
		return mergedTerms
	}
	if prop == "CustomSortOrder" {
		if val == nil {
			return nil
		}
		sortedTerms := strings.Split(val.(string), ":")
		for i, term := range sortedTerms {
			sortedTerms[i] = strings.ReplaceAll(strings.ReplaceAll(term, "/Guid(", ""), ")/", "")
		}
		return sortedTerms
	}
	return val
}
