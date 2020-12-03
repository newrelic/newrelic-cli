package install

import "context"

type recipeExecutor interface {
	execute(context.Context, discoveryManifest, recipe) error
}
