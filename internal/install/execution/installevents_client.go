package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type InstalleventsClient interface {
	CreateRecipeEvent(int, installevents.RecipeStatus) (*installevents.RecipeEvent, error)
}
