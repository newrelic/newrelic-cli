package install

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type regexProcessFilterer struct {
	recipeFetcher recipeFetcher
}

func newRegexProcessFilterer(r recipeFetcher) *regexProcessFilterer {
	f := regexProcessFilterer{
		recipeFetcher: r,
	}

	return &f
}

func (f *regexProcessFilterer) filter(ctx context.Context, processes []genericProcess) ([]genericProcess, error) {
	recipes, err := f.recipeFetcher.fetchRecipes(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve process filter criteria: %s", err)
	}

	matches := []genericProcess{}
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

func match(r recipe, process genericProcess) bool {
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
