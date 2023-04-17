package lists

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudquery/plugin-sdk/v2/plugins/source"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
	"github.com/thoas/go-funk"
)

func (l *Lists) Sync(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table) error {
	opts := l.TablesMap[table.Name]
	logger := l.logger.With().Str("table", table.Name).Logger()

	logger.Debug().Strs("cols", opts.Spec.Select).Msg("selecting columns from list")

	top := 5000
	if opts.Spec.Top > 0 && opts.Spec.Top < 5000 {
		top = opts.Spec.Top
	}

	list := l.sp.Web().GetList(opts.URI)
	items, err := list.Items().
		Select(strings.Join(opts.Spec.Select, ",")).
		Expand(strings.Join(opts.Spec.Expand, ",")).
		Filter(opts.Spec.Filter).
		Top(top).GetPaged()

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

			colVals := make([]any, len(table.Columns))

			for i, col := range table.Columns {
				prop := col.Description
				colVals[i] = util.GetRespValByProp(itemMap, prop)
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
		if col.Type == schema.TypeString {
			if values[i] != nil {
				values[i] = fmt.Sprintf("%v", values[i])
			}
		}
		if err := resource.Set(col.Name, values[i]); err != nil {
			return nil, err
		}
	}
	return resource, nil
}
