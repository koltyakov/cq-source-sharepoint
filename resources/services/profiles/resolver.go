package profiles

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/gosip/api"
)

func (u *Profiles) Resolver(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	rowLimit := 500
	startRow := 0

	data, err := searchUsers(u.sp, startRow, rowLimit)

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
