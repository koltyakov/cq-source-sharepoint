package search

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/koltyakov/gosip/api"
	"github.com/thoas/go-funk"
)

func (s *Search) Sync(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table) error {
	opts := s.TablesMap[table.Name]

	rowLimit := 500
	startRow := 0

	data, err := searchData(s.sp, opts.Spec, startRow, rowLimit)

	for {
		if err != nil {
			metrics.Errors++
			return fmt.Errorf("failed to get items: %w", err)
		}

		rows := data.Data().PrimaryQueryResult.RelevantResults.Table.Rows

		for _, row := range rows {
			ks := funk.Keys(row).([]string)
			sort.Strings(ks)

			colVals := make([]any, len(table.Columns))

			for i, col := range table.Columns {
				prop := col.Description
				colVals[i] = getSearchCellValue(row, prop)
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

		if len(rows) < rowLimit {
			break
		}
		startRow += rowLimit
		data, err = searchData(s.sp, opts.Spec, startRow, rowLimit)
	}

	return nil
}

func searchData(sp *api.SP, spec Spec, startRow int, rowLimit int) (api.SearchResp, error) {
	return sp.Search().PostQuery(&api.SearchQuery{
		QueryText:        spec.QueryText,
		SourceID:         spec.SourceID,
		SelectProperties: spec.SelectProperties,
		TrimDuplicates:   spec.TrimDuplicates,
		StartRow:         startRow,
		RowLimit:         rowLimit,
	})
}

func getSearchCellValue(row *struct {
	Cells []*api.TypedKeyValue `json:"Cells"`
}, prop string) any {
	for _, cell := range row.Cells {
		if cell.Key == prop {
			if prop == "UniqueId" {
				return strings.ReplaceAll(strings.ReplaceAll(cell.Value, "{", ""), "}", "")
			}
			return cell.Value
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
