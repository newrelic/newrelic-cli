package recipes

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RegexProcessMatchFinder struct{}

func NewRegexProcessMatchFinder() *RegexProcessMatchFinder {
	f := RegexProcessMatchFinder{}

	return &f
}

func (f *RegexProcessMatchFinder) FindMatches(ctx context.Context, processes []types.GenericProcess, recipe types.OpenInstallationRecipe) ([]types.MatchedProcess, error) {

	matches := []types.MatchedProcess{}
	for _, p := range processes {
		m, err := f.findMatches(recipe, p)
		if err != nil {
			return nil, err
		}

		matches = append(matches, m...)
	}

	if len(matches) > 0 {
		log.Debugf("Finished matching recipe %s to running processes, found %d matches.", recipe.Name, len(matches))
	}
	return matches, nil
}

func (f *RegexProcessMatchFinder) FindMatchesMultiple(ctx context.Context, processes []types.GenericProcess, recipes []types.OpenInstallationRecipe) ([]types.MatchedProcess, error) {
	matches := []types.MatchedProcess{}
	log.Debugf("Filtering recipes with %d processes...", len(processes))

	for _, r := range recipes {
		m, err := f.FindMatches(ctx, processes, r)
		if err != nil {
			return nil, err
		}

		matches = append(matches, m...)
	}

	if len(matches) > 0 {
		log.Debugf("Filtering recipes with processes done, found %d matches.", len(matches))
	}
	return matches, nil
}

func (f *RegexProcessMatchFinder) findMatches(r types.OpenInstallationRecipe, process types.GenericProcess) ([]types.MatchedProcess, error) {
	matches := []types.MatchedProcess{}
	for _, pattern := range r.ProcessMatch {
		cmd, err := process.Cmd()
		if err != nil {
			return nil, err
		}

		matched, err := regexp.Match(pattern, []byte(cmd))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, cmd)
			continue
		}

		if matched {
			mp, err := makeMatchedProcess(process)
			if err != nil {
				return nil, err
			}
			mp.MatchingPattern = pattern
			mp.MatchingRecipe = r
			log.Debugf("Process matching pattern %s with %s for recipe %s.", pattern, cmd, r.DisplayName)

			matches = append(matches, *mp)
		}
	}

	return matches, nil
}

func makeMatchedProcess(p types.GenericProcess) (*types.MatchedProcess, error) {
	cmdLine, err := p.Cmd()
	if err != nil {
		return nil, err
	}

	if cmdLine == "" {
		return nil, fmt.Errorf("empty command for pid %d", p.PID())
	}

	return &types.MatchedProcess{
		GenericProcess: p,
	}, nil
}
