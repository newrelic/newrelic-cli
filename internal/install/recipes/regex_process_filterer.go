package recipes

import (
	"context"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RegexProcessFilterer struct {
}

func NewRegexProcessFilterer() *RegexProcessFilterer {
	f := RegexProcessFilterer{}

	return &f
}

func (f *RegexProcessFilterer) Filter(ctx context.Context, processes []types.GenericProcess, recipes []types.OpenInstallationRecipe) []types.MatchedProcess {
	matchedProcesses := getMatchedProcesses(processes)
	log.Debugf("Filtering recipes with %d processes...", len(matchedProcesses))

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
	return matches
}

func match(r types.OpenInstallationRecipe, matchedProcess *types.MatchedProcess) bool {
	for _, pattern := range r.ProcessMatch {
		matched, err := regexp.Match(pattern, []byte(matchedProcess.Command))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, matchedProcess.Command)
			continue
		}

		if matched {
			matchedProcess.MatchingPattern = pattern
			matchedProcess.MatchingRecipe = r
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
