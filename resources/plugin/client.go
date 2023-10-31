package plugin

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/cq-source-sharepoint/resources/auth"
	"github.com/rs/zerolog"
)

type Client struct {
	logger    zerolog.Logger
	spec      Spec
	tables    schema.Tables
	scheduler *scheduler.Scheduler

	options plugin.NewClientOptions

	plugin.UnimplementedDestination
}

func (*Client) ID() string {
	return Name
}
func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	if c.options.NoConnection {
		return fmt.Errorf("no connection")
	}

	tt, err := c.Tables(ctx, plugin.TableOptions{
		Tables:              options.Tables,
		SkipTables:          options.SkipTables,
		SkipDependentTables: options.SkipDependentTables,
	})

	if err != nil {
		return err
	}

	return c.scheduler.Sync(ctx, c, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	if c.options.NoConnection {
		return schema.Tables{}, nil
	}

	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (*Client) Close(context.Context) error {
	// ToDo: Add your client cleanup here
	return nil
}

func NewClient(_ context.Context, logger zerolog.Logger, cnfg []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	logger = logger.With().Str("plugin", "sharepoint").Logger()

	if opts.NoConnection {
		// no spec could be present
		return &Client{
			logger:  logger,
			options: opts,
		}, nil
	}

	spec, err := getSpec(cnfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}

	sp, err := auth.GetSP(spec.Auth)
	if err != nil {
		return nil, err
	}

	if _, err := sp.Web().Select("Title").Get(); err != nil {
		return nil, fmt.Errorf("failed to connect to SharePoint: %w", err)
	}

	tables, err := spec.getTables(sp, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tables: %w", err)
	}

	return &Client{
		logger:    logger,
		spec:      *spec,
		tables:    tables,
		scheduler: scheduler.NewScheduler(scheduler.WithLogger(logger)),
		options:   opts,
	}, nil
}
