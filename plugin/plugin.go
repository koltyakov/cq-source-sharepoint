package plugin

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/cq-source-sharepoint/client"
)

var (
	Name    = "cq-source-sharepoint"
	Kind    = "source"
	Team    = "koltyakov"
	Version = "development"
)

func Plugin() *plugin.Plugin {
	return plugin.NewPlugin(
		Name,
		Version,
		client.NewClient,
		plugin.WithKind(Kind),
		plugin.WithTeam(Team),
		// source.WithDynamicTableOption(getDynamicTables),
		// source.WithUnmanaged(),
		// source.WithNoInternalColumns(),
	)
}

func getDynamicTables(ctx context.Context, c schema.ClientMeta) (schema.Tables, error) {
	cl := c.(*client.Client)
	return cl.Tables, nil
}
