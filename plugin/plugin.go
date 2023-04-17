package plugin

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v2/plugins/source"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/koltyakov/cq-source-sharepoint/client"
)

var (
	Version = "development"
)

func Plugin() *source.Plugin {
	return source.NewPlugin(
		"sharepoint",
		Version,
		nil,
		client.NewClient,
		source.WithDynamicTableOption(getDynamicTables),
		source.WithUnmanaged(),
		source.WithNoInternalColumns(),
	)
}

func getDynamicTables(ctx context.Context, c schema.ClientMeta) (schema.Tables, error) {
	cl := c.(*client.Client)
	return cl.Tables, nil
}
