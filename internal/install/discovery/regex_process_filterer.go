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

func (f *RegexProcessFilterer) filter(ctx context.Context, processes []types.GenericProcess) ([]types.ProcessInfoWrap, error) {
	processesInfo := getProcessesInfo(processes)
	log.Debugf("Filtering recipes with %d processes...", len(processesInfo))
	for _, p := range processesInfo {
		log.Debugf("Match using processInfo: %s", p.Info)
	}

	recipes, err := f.recipeFetcher.FetchRecipes(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve process filter criteria: %s", err)
	}
	for _, r := range recipes {
		log.Debugf("Match using recipe DisplayName: %s RecipeProcessMatch: %s", r.DisplayName, r.ProcessMatch)
	}

	matches := []types.ProcessInfoWrap{}
	for _, p := range processesInfo {
		isMatch := false
		for _, r := range recipes {
			isMatch = isMatch || match(r, &p)
		}

		if isMatch {
			matches = append(matches, p)
		}
	}

	log.Debugf("Filtering recipes with processes done, found %d matches.", len(matches))
	return matches, nil
}

func match(r types.Recipe, processInfo *types.ProcessInfoWrap) bool {
	for _, pattern := range r.ProcessMatch {
		matched, err := regexp.Match(pattern, []byte(processInfo.Info))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, processInfo.Info)
			continue
		}

		if matched {
			processInfo.MatchingPattern = pattern
			log.Debugf("Process matching pattern %s with %s for recipe %s.", pattern, processInfo.Info, r.DisplayName)
			return matched
		}
	}

	return false
}

func getProcessesInfo(processes []types.GenericProcess) []types.ProcessInfoWrap {
	processesInfo := []types.ProcessInfoWrap{}
	for _, p := range processes {
		cmdLine, err := p.Cmdline()
		if err == nil {
			if len(cmdLine) > 0 {
				processInfo := types.ProcessInfoWrap{
					Info:    cmdLine,
					Process: p,
				}
				processesInfo = append(processesInfo, processInfo)
			}
		}
	}
	return processesInfo
}
