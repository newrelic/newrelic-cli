package execution

import "github.com/newrelic/newrelic-client-go/v2/pkg/installevents"

type InstallEventsClient interface {
	InstallationCreateRecipeEvent(int, installevents.InstallationRecipeStatus) (*installevents.InstallationRecipeEvent, error)
	InstallationCreateInstallStatus(int, installevents.InstallationInstallStatusInput) (*installevents.InstallationInstallStatus, error)
}
