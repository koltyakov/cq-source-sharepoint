package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/gosip/api"
)

type ResolverClosure = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error

func (s *Search) Resolver(spec Spec, table *schema.Table) ResolverClosure {
	return func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		rowLimit := 500
		startRow := 0

		data, err := searchData(s.sp, spec, startRow, rowLimit)

		for {
			if err != nil {
				return fmt.Errorf("failed to get items: %w", err)
			}

			rows := data.Data().PrimaryQueryResult.RelevantResults.Table.Rows

			select {
			case <-ctx.Done():
				return ctx.Err()
			case res <- rows:
			}

			if len(rows) < rowLimit {
				break
			}
			startRow += rowLimit
			data, err = searchData(s.sp, spec, startRow, rowLimit)
		}

		return nil
	}
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
