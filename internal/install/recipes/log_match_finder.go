package recipes

import (
	"context"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type LogMatchFinder struct{}

type LogMatchFinderDefinition interface {
	GetPaths(context.Context, []*types.OpenInstallationRecipe) []types.OpenInstallationLogMatch
}

func NewLogMatchFinder() LogMatchFinderDefinition {
	f := LogMatchFinder{}

	return &f
}

func (f *LogMatchFinder) GetPaths(ctx context.Context, recipes []*types.OpenInstallationRecipe) []types.OpenInstallationLogMatch {
	fileMatches := []types.OpenInstallationLogMatch{}

	for _, r := range recipes {
		for _, l := range r.LogMatch {
			match, _ := matchLogFilesFromRecipe(l)
			if match {
				fileMatches = append(fileMatches, l)
			}
		}
	}

	return fileMatches
}

func matchLogFilesFromRecipe(matcher types.OpenInstallationLogMatch) (bool, []string) {
	matches, err := filepath.Glob(matcher.File)
	if err != nil {
		log.Errorf("error matching logfiles: %s", err)
		return false, nil
	}

	if len(matches) > 0 {
		return true, matches
	}

	return false, nil
}
