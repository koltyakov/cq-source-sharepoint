package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/addin"
	"github.com/rs/zerolog"
)

type Client struct {
	Logger zerolog.Logger
	Tables schema.Tables
	SP 	 *api.SP
}

func (c *Client) ID() string {
	// TODO: Change to either your plugin name or a unique dynamic identifier
	return "ID"
}

func New(ctx context.Context, logger zerolog.Logger, s specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var pluginSpec Spec

	if err := s.UnmarshalSpec(&pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}
	if err := pluginSpec.Validate(); err != nil {
		return nil, err
	}

	auth := &strategy.AuthCnfg{
		SiteURL:     pluginSpec.SiteURL,
		ClientID:    pluginSpec.ClientID,
		ClientSecret: pluginSpec.ClientSecret,
	}
	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)
	lists, err := sp.Web().Lists().Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}
	tables := make(schema.Tables, len(lists.Data()))
	normalizedNames := make(map[string]bool)
	for i, list := range lists.Data(){
		name := normalizeList(list.Data().Title)
		if _, ok := normalizedNames[name]; ok {
			logger.Warn().Msgf("List %s has been normalized to %s, but another list has already been normalized to that name. skipping %s", list.Data().Title, name, list.Data().Title)
			continue
		}
		normalizedNames[name] = true
		table := &schema.Table{
			Name: "sharepoint_" + name,
			Description: list.Data().Title,
		}
		fields, err := sp.Web().GetList("Lists/" + list.Data().Title).Fields().Get()
		if err != nil {
			if IsNotFound(err) {
				continue
			}
			return nil, fmt.Errorf("failed to get fields: %w", err)
		}

		for _, field := range fields.Data() {
			fmt.Println("key:" + field.Data().Title)
			fmt.Println("field:" + field.Data().TypeAsString)
		}
		fmt.Println(table.Name)
		tables[i] = table
	}

	return &Client{
		Logger: logger,
		Tables: tables,
		SP: sp,
	}, nil
}


func normalizeList(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}