package aimonitoring

import (
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

const maxEntitySearchPages = 1000

type entitySearchResponse[T any] struct {
	Actor struct {
		EntityManagement struct {
			EntitySearch struct {
				Entities   []T    `json:"entities"`
				NextCursor string `json:"nextCursor"`
			} `json:"entitySearch"`
		} `json:"entityManagement"`
	} `json:"actor"`
}

func fetchAllEntities[T any](query string) ([]T, error) {
	var all []T

	var cursor *string
	for range maxEntitySearchPages {
		var resp entitySearchResponse[T]
		vars := map[string]interface{}{"cursor": cursor}
		if err := client.NRClient.NerdGraph.QueryWithResponseAndContext(utils.SignalCtx, query, vars, &resp); err != nil {
			return nil, err
		}

		search := resp.Actor.EntityManagement.EntitySearch
		all = append(all, search.Entities...)

		if search.NextCursor == "" {
			break
		}
		nextCursor := search.NextCursor
		cursor = &nextCursor
	}

	return all, nil
}
