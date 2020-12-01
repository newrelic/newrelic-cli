package install

import "context"

type processFilterer interface {
	filter(context.Context, []genericProcess) ([]genericProcess, error)
}
