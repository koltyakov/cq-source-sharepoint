package ct

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/koltyakov/gosip/api"
)

type ResolverClosure = func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error

func (c *ContentTypesRollup) Resolver(contentTypeID string, spec Spec, table *schema.Table) ResolverClosure {
	return func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
		logger := c.logger.With().Str("table", table.Name).Logger()

		logger.Debug().Msgf("getting webs for %s", table.Name)

		webUrls, err := c.getWebs(c.sp.ToURL())
		if err != nil {
			return err
		}
		logger.Debug().Msgf("webs found: %v", webUrls)

		// Iterate over all webs
		for _, webURL := range webUrls {
			logger.Debug().Msgf("getting lists for %s", webURL)
			lists, err := c.getLists(webURL, contentTypeID)
			if err != nil {
				return err
			}
			logger.Debug().Msgf("lists with content type: %v", lists)

			// Iterate over all lists
			for _, listID := range lists {
				c.logger.Debug().Msgf("list sync: %s", listID)
				if err := c.syncList(ctx, webURL, listID, contentTypeID, res, spec); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func (c *ContentTypesRollup) getWeb(webURL string) *api.Web {
	return c.sp.Web().FromURL(fmt.Sprintf("%s/_api/Web", webURL))
}

func (c *ContentTypesRollup) getWebs(webURL string) ([]string, error) {
	web := c.getWeb(webURL)

	resp, err := web.Webs().Select("Url,Webs/Url").Expand("Webs").Top(5000).Get()
	if err != nil {
		return nil, err
	}

	var webs []struct {
		URL  string `json:"Url"`
		Webs []struct {
			URL string `json:"Url"`
		} `json:"Webs"`
	}

	if err := json.Unmarshal(resp.Normalized(), &webs); err != nil {
		return nil, err
	}

	webURLs := []string{webURL}
	for _, web := range webs {
		webURLs = append(webURLs, web.URL)
		for _, subWeb := range web.Webs {
			subWebs, err := c.getWebs(subWeb.URL)
			if err != nil {
				return nil, err
			}
			webURLs = append(webURLs, subWebs...)
		}
	}

	return webURLs, nil
}

func (c *ContentTypesRollup) getLists(webURL string, ctID string) ([]string, error) {
	web := c.getWeb(webURL)

	resp, err := web.Lists().
		Select("Id,ContentTypes/StringId").
		Filter("AllowContentTypes eq true and Hidden eq false").
		Expand("ContentTypes").
		Top(5000).Get()

	if err != nil {
		return nil, err
	}

	var lists []struct {
		ID           string `json:"Id"`
		ContentTypes []struct {
			StringId string `json:"StringId"`
		} `json:"ContentTypes"`
	}
	if err := json.Unmarshal(resp.Normalized(), &lists); err != nil {
		return nil, err
	}

	var listIds []string
	for _, list := range lists {
		for _, ct := range list.ContentTypes {
			if strings.HasPrefix(ct.StringId, ctID) {
				listIds = append(listIds, list.ID)
			}
		}
	}

	return listIds, nil
}

func (c *ContentTypesRollup) syncList(ctx context.Context, webURL string, listID string, ctID string, res chan<- any, spec Spec) error {
	web := c.getWeb(webURL)
	list := web.Lists().GetByID(listID)

	// Content type is not applied as filter in query to support lists of any size
	// it is used to filter results in memory after getting responses
	items, err := list.Items().
		Select(strings.Join(append(spec.Select, "ContentTypeId"), ",")).
		Expand(strings.Join(spec.Expand, ",")).
		Top(5000).
		GetPaged()

	for {
		if err != nil {
			return fmt.Errorf("failed to get items: %w", err)
		}

		var itemList []map[string]any
		if err := json.Unmarshal(items.Items.Normalized(), &itemList); err != nil {
			return err
		}

		for _, itemMap := range itemList {
			// Filter by content type ID (skip items which content type is doesn't strart with base content type ID)
			if !strings.HasPrefix(itemMap["ContentTypeId"].(string), ctID) {
				continue
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case res <- itemMap:
			}
		}

		if !items.HasNextPage() {
			break
		}
		items, err = items.GetNextPage()
	}

	return nil
}
