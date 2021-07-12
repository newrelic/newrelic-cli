package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type InstalleventsClient interface {
	InstallationCreateRecipeEvent(int, installevents.InstallationRecipeStatus) (*installevents.InstallationRecipeEvent, error)
}
