package mmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/cq-source-sharepoint/internal/util"
)

type ResolverClosure = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error

func (m *MMD) Resolver(terSetID string, spec Spec, table *schema.Table) ResolverClosure {
	return func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		taxonomy := m.sp.Taxonomy()
		terms, err := taxonomy.Stores().Default().Sets().GetByID(terSetID).GetAllTerms()
		if err != nil {
			return fmt.Errorf("failed to get items: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case res <- terms:
		}

		return nil
	}
}

func getRespValByProp(resp map[string]any, prop string) any {
	val := util.GetRespValByProp(resp, prop)
	if prop == "Id" {
		return strings.ReplaceAll(strings.ReplaceAll(val.(string), "/Guid(", ""), ")/", "")
	}
	if prop == "CreatedDate" || prop == "LastModifiedDate" {
		dateStr := strings.ReplaceAll(strings.ReplaceAll(val.(string), "/Date(", ""), ")/", "")
		dateInt, _ := strconv.ParseInt(dateStr, 10, 64)
		return time.UnixMilli(dateInt)
	}
	if prop == "PathOfTerm" {
		return strings.Split(val.(string), ";")
	}
	if prop == "MergedTermIds" {
		mergedTerms := val.([]any)
		for i, term := range mergedTerms {
			mergedTerms[i] = strings.ReplaceAll(strings.ReplaceAll(term.(string), "/Guid(", ""), ")/", "")
		}
		return mergedTerms
	}
	if prop == "CustomSortOrder" {
		if val == nil {
			return nil
		}
		sortedTerms := strings.Split(val.(string), ":")
		for i, term := range sortedTerms {
			sortedTerms[i] = strings.ReplaceAll(strings.ReplaceAll(term, "/Guid(", ""), ")/", "")
		}
		return sortedTerms
	}
	return val
}
