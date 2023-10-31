package plugin

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/ct"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/lists"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/mmd"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/profiles"
	"github.com/koltyakov/cq-source-sharepoint/resources/services/search"
	"github.com/koltyakov/gosip/api"
	"github.com/rs/zerolog"
)

func (s *Spec) getTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	tables := schema.Tables{}

	// Tables from lists config
	listTables, err := s.getListsTables(sp, logger)
	if err != nil {
		return nil, err
	}
	tables = append(tables, listTables...)

	// Tables from mmd config
	mmdTables, err := s.getMMDTables(sp, logger)
	if err != nil {
		return nil, err
	}
	tables = append(tables, mmdTables...)

	// Tables from profiles config
	profileTables, err := s.getProfileTables(sp, logger)
	if err != nil {
		return nil, err
	}
	tables = append(tables, profileTables...)

	// Tables from search config
	searchTables, err := s.getSearchTables(sp, logger)
	if err != nil {
		return nil, err
	}
	tables = append(tables, searchTables...)

	// Tables from content types config
	ctTables, err := s.getContentTypeTables(sp, logger)
	if err != nil {
		return nil, err
	}
	tables = append(tables, ctTables...)

	if err := transformers.TransformTables(tables); err != nil {
		return nil, err
	}

	for _, table := range tables {
		schema.AddCqIDs(table)
	}

	return tables, nil
}

func (s *Spec) getListsTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	tables := make(schema.Tables, 0, len(s.Lists))
	l := lists.NewLists(sp, logger)
	for uri, spec := range s.Lists {
		table, err := l.GetDestTable(uri, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get list '%s': %w", uri, err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (s *Spec) getMMDTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	tables := make(schema.Tables, 0, len(s.MMD))
	m := mmd.NewMMD(sp, logger)
	for id, spec := range s.MMD {
		table, err := m.GetDestTable(id, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get term set '%s': %w", id, err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (s *Spec) getProfileTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	if !s.Profiles.Enabled {
		return nil, nil
	}

	table, err := profiles.NewProfiles(sp, logger).GetDestTable(s.Profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %w", err)
	}
	return schema.Tables{table}, nil
}

func (s *Spec) getSearchTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	tables := make(schema.Tables, 0, len(s.Search))
	srch := search.NewSearch(sp, logger)
	for name, spec := range s.Search {
		table, err := srch.GetDestTable(name, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get search '%s': %w", name, err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (s *Spec) getContentTypeTables(sp *api.SP, logger zerolog.Logger) (schema.Tables, error) {
	tables := make(schema.Tables, 0, len(s.ContentTypes))
	c := ct.NewContentTypesRollup(sp, logger)
	for name, spec := range s.ContentTypes {
		table, err := c.GetDestTable(name, spec)
		if err != nil {
			return nil, fmt.Errorf("failed to get content type '%s': %w", name, err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}
