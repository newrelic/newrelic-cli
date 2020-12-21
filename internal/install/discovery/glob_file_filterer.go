package discovery

import (
	"context"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// GlobFileFilterer is an implementation of the FileFilterer interface that uses
// glob-based filesystem searches to locate the existence of files.
type GlobFileFilterer struct{}

// NewGlobFileFilterer returns a new instance of GlobFileFilterer.
func NewGlobFileFilterer() *GlobFileFilterer {
	f := GlobFileFilterer{}

	return &f
}

// Filter uses the patterns provided in the passed recipe to return matches based
// on which files exist in the underlying file system.
func (f *GlobFileFilterer) Filter(ctx context.Context, recipes []types.Recipe) ([]types.LogMatch, error) {
	fileMatches := []types.LogMatch{}
	for _, r := range recipes {
		for _, l := range r.LogMatch {
			match, _ := matchLogFilesFromRecipe(l)
			if match {
				fileMatches = append(fileMatches, l)
			}
		}
	}

	return fileMatches, nil
}

func matchLogFilesFromRecipe(matcher types.LogMatch) (bool, []string) {
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
