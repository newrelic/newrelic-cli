package aimonitoring

import (
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

const maxEntitySearchPages = 1000

// fetchAllEntities pages through entityManagement.entitySearch and returns every entity of type T.
func fetchAllEntities[T aimonitoring.EntityManagementEntityInterface](query string) ([]T, error) {
	var all []T

	var cursor *string
	for range maxEntitySearchPages {
		result, err := client.NRClient.AIMonitoring.GetEntitySearchWithContext(utils.SignalCtx, cursor, query)
		if err != nil {
			return nil, err
		}

		for _, e := range result.Entities {
			if entity, ok := e.(T); ok {
				all = append(all, entity)
			}
		}

		if result.NextCursor == "" {
			break
		}
		cursor = &result.NextCursor
	}

	return all, nil
}
