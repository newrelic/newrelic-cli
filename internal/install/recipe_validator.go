package install

import "context"

type recipeValidator interface {
	validate(context.Context, discoveryManifest, recipe) (ok bool, entityGUID string, err error)
}
