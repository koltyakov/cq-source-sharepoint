package lists

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type ResolverClosure = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error

func (l *Lists) Resolver(listURI string, spec Spec, table *schema.Table) ResolverClosure {
	return func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		logger := l.logger.With().Str("table", table.Name).Logger()

		logger.Debug().Strs("cols", spec.Select).Msg("selecting columns from list")

		top := 5000
		if spec.Top > 0 && spec.Top < 5000 {
			top = spec.Top
		}

		list := l.sp.Web().GetList(listURI)
		items, err := list.Items().
			Select(strings.Join(spec.Select, ",")).
			Expand(strings.Join(spec.Expand, ",")).
			Filter(spec.Filter).
			Top(top).GetPaged()

		for {
			if err != nil {
				return fmt.Errorf("failed to get items: %w", err)
			}

			var itemList []map[string]any
			if err := json.Unmarshal(items.Items.Normalized(), &itemList); err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case res <- itemList:
			}

			if !items.HasNextPage() {
				break
			}
			items, err = items.GetNextPage()
		}

		return nil
	}
}
