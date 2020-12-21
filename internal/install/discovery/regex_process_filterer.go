package discovery

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RegexProcessFilterer struct {
	recipeFetcher recipes.RecipeFetcher
}

func NewRegexProcessFilterer(r recipes.RecipeFetcher) *RegexProcessFilterer {
	f := RegexProcessFilterer{
		recipeFetcher: r,
	}

	return &f
}

func (f *RegexProcessFilterer) filter(ctx context.Context, processes []types.GenericProcess) ([]types.GenericProcess, error) {
	recipes, err := f.recipeFetcher.FetchRecipes(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve process filter criteria: %s", err)
	}

	matches := []types.GenericProcess{}
	for _, p := range processes {
		isMatch := false
		for _, r := range recipes {
			isMatch = isMatch || match(r, p)
		}

		if isMatch {
			matches = append(matches, p)
		}
	}

	return matches, nil
}

func match(r types.Recipe, process types.GenericProcess) bool {
	for _, pattern := range r.ProcessMatch {
		name, err := process.Name()
		if err != nil {
			log.Debugf("could not retrieve process name for PID %d", process.PID())
			continue
		}

		matched, err := regexp.Match(pattern, []byte(name))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, name)
			continue
		}

		if matched {
			return matched
		}
	}

	return false
}
