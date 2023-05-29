package profiles

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudquery/plugin-sdk/v3/plugins/source"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/koltyakov/gosip/api"
	"github.com/thoas/go-funk"
)

func (u *Profiles) Sync(ctx context.Context, metrics *source.TableClientMetrics, res chan<- *schema.Resource, table *schema.Table) error {
	rowLimit := 500
	startRow := 0

	data, err := searchUsers(u.sp, startRow, rowLimit)

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
		data, err = searchUsers(u.sp, startRow, rowLimit)
	}

	return nil
}

func searchUsers(sp *api.SP, startRow int, rowLimit int) (api.SearchResp, error) {
	return sp.Search().PostQuery(&api.SearchQuery{
		QueryText:          "*",
		SourceID:           "b09a7990-05ea-4af9-81ef-edfab16c4e31",
		SelectProperties:   userProps,
		TrimDuplicates:     false,
		EnableInterleaving: true,
		StartRow:           startRow,
		RowLimit:           rowLimit,
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
