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

func (f *RegexProcessFilterer) filter(ctx context.Context, processes []types.GenericProcess, manifest types.DiscoveryManifest) ([]types.MatchedProcess, error) {
	matchedProcesses := getMatchedProcesses(processes)
	log.Debugf("Filtering recipes with %d processes...", len(matchedProcesses))

	recipes, err := f.recipeFetcher.FetchRecipes(ctx, &manifest)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve process filter criteria: %s", err)
	}

	for _, r := range recipes {
		log.Tracef("Match using recipe DisplayName: %s RecipeProcessMatch: %s", r.DisplayName, r.ProcessMatch)
	}

	matches := []types.MatchedProcess{}
	for _, p := range matchedProcesses {
		log.Tracef("Match using process command: %s", p.Command)
		for _, r := range recipes {
			if match(r, &p) {
				matches = append(matches, p)
			}
		}
	}

	log.Debugf("Filtering recipes with processes done, found %d matches.", len(matches))
	return matches, nil
}

func match(r types.Recipe, matchedProcess *types.MatchedProcess) bool {
	for _, pattern := range r.ProcessMatch {
		matched, err := regexp.Match(pattern, []byte(matchedProcess.Command))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, matchedProcess.Command)
			continue
		}

		if matched {
			matchedProcess.MatchingPattern = pattern
			log.Debugf("Process matching pattern %s with %s for recipe %s.", pattern, matchedProcess.Command, r.DisplayName)
			return matched
		}
	}

	return false
}

func getMatchedProcesses(processes []types.GenericProcess) []types.MatchedProcess {
	matchedProcesses := []types.MatchedProcess{}
	for _, p := range processes {
		cmdLine, err := p.Cmdline()
		if err != nil {
			continue
		}
		if cmdLine != "" {
			matchedProcess := types.MatchedProcess{
				Command: cmdLine,
				Process: p,
			}
			matchedProcesses = append(matchedProcesses, matchedProcess)
		}
	}
	return matchedProcesses
}
