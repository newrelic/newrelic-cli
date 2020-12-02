package install

import "context"

type recipeValidator interface {
	validate(context.Context, recipe) (bool, error)
}
