package install

type executionStatusReporter interface {
	reportRecipeFailed(event recipeStatusEvent) error
	reportRecipeInstalled(event recipeStatusEvent) error
	reportRecipesAvailable(recipes []recipe) error
	reportComplete() error
}

type recipeStatusEvent struct {
	recipe     recipe
	msg        string
	entityGUID string
}
