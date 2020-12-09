package install

import "context"

type fileFilterer interface {
	filter(context.Context, []recipe) ([]logMatch, error)
}
