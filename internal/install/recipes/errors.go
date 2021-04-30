package recipes

import "github.com/pkg/errors"

// ErrRecipeNotFound is used when a recipe is requested by name, but does not exist for the given constraint.
var ErrRecipeNotFound = errors.New("recipe not found")
