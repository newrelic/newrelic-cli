package recipes

import (
	"context"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RegexProcessMatchFinder struct{}

func NewRegexProcessMatchFinder() *RegexProcessMatchFinder {
	f := RegexProcessMatchFinder{}

	return &f
}

func (f *RegexProcessMatchFinder) FindMatches(ctx context.Context, processes []types.GenericProcess, recipe types.OpenInstallationRecipe) []types.MatchedProcess {

	matches := []types.MatchedProcess{}
	for _, p := range processes {
		m := f.findMatches(recipe, p)

		matches = append(matches, m...)
	}

	if len(matches) > 0 {
		log.Debugf("Finished matching recipe %s to running processes, found %d matches.", recipe.Name, len(matches))
	}
	return matches
}

func (f *RegexProcessMatchFinder) FindMatchesMultiple(ctx context.Context, processes []types.GenericProcess, recipes []types.OpenInstallationRecipe) []types.MatchedProcess {
	matches := []types.MatchedProcess{}
	log.Debugf("Filtering recipes with %d processes...", len(processes))

	for _, r := range recipes {
		m := f.FindMatches(ctx, processes, r)

		matches = append(matches, m...)
	}

	if len(matches) > 0 {
		log.Debugf("Filtering recipes with processes done, found %d matches.", len(matches))
	}
	return matches
}

func (f *RegexProcessMatchFinder) findMatches(r types.OpenInstallationRecipe, process types.GenericProcess) []types.MatchedProcess {
	matches := []types.MatchedProcess{}
	var newrelicInstallRegex = regexp.MustCompile(`(?i).newrelic(|\.exe['"]?) install.`)
	for _, pattern := range r.ProcessMatch {
		cmd, err := process.Cmd()
		if err != nil {
			// Process no longer exist, skip
			continue
		}

		if len(newrelicInstallRegex.FindString(cmd)) > 0 {
			continue
		}

		matched, err := regexp.Match(pattern, []byte(cmd))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s, %s", pattern, cmd, err)
			continue
		}

		if matched {
			mp := &types.MatchedProcess{}
			mp.GenericProcess = process
			mp.MatchingPattern = pattern
			mp.MatchingRecipe = r
			log.Debugf("Process matching pattern %s with %s for recipe %s.", pattern, cmd, r.DisplayName)

			matches = append(matches, *mp)
		}
	}

	return matches
}
