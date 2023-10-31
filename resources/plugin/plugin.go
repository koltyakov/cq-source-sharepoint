package plugin

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

var (
	Name    = "sharepoint"
	Kind    = "source"
	Team    = "koltyakov"
	Version = "development"
)

type Plugin struct {
	logger    zerolog.Logger
	spec      Spec
	tables    schema.Tables
	scheduler *scheduler.Scheduler

	client *Client

	plugin.UnimplementedDestination
}

func NewPlugin() *plugin.Plugin {
	return plugin.NewPlugin(Name, Version, NewClient, plugin.WithKind(Kind), plugin.WithTeam(Team))
}

func (*Plugin) ID() string {
	return Name
}

func (p *Plugin) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := p.Tables(ctx, plugin.TableOptions{
		Tables:              options.Tables,
		SkipTables:          options.SkipTables,
		SkipDependentTables: options.SkipDependentTables,
	})

	if err != nil {
		return err
	}

	return p.scheduler.Sync(ctx, p, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (p *Plugin) Tables(ctx context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := p.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (*Plugin) Close(ctx context.Context) error {
	// ToDo: Add your client cleanup here
	return nil
}
