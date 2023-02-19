package client

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
)

func (c *Client) Sync(ctx context.Context, metrics *source.Metrics, res chan<- *schema.Resource) error {
	for _, table := range c.Tables {
		list := c.SP.Web().GetList("Lists/" + table.Description)
		items, err := list.Items().GetAll()
		if err != nil {
			if IsNotFound(err) {
				continue
			}
			return fmt.Errorf("failed to get items: %w", err)
		}
		for _, item := range items {
			fmt.Println(item.ToMap())
		}
	}

	return nil
}