package install

import "context"

type recipeValidator interface {
	validate(context.Context, discoveryManifest, recipe) (bool, error)
}
