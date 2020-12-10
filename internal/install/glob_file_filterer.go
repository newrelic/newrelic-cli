package install

import (
	"context"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type globFileFilterer struct {
	// recipeFetcher recipeFetcher
}

func newGlobFileFilterer() *globFileFilterer {
	f := globFileFilterer{}

	return &f
}

func (f *globFileFilterer) filter(ctx context.Context, recipes []recipe) ([]logMatch, error) {
	fileMatches := []logMatch{}
	for _, r := range recipes {
		for _, l := range r.LogMatch {

			match, _ := matchLogFilesFromRecipe(l)
			if match {
				// if userAcceptLogFiles(files) {
				fileMatches = append(fileMatches, l)
				// }
			}

		}
	}

	return fileMatches, nil
}

func matchLogFilesFromRecipe(matcher logMatch) (bool, []string) {
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

// func userAcceptLogFiles(files []string) bool {
// 	msg := fmt.Sprintf("The following log files have been found: %s\nDo you want to watch them? [Yes/No]", strings.Join(files, ", "))
//
// 	prompt := promptui.Select{
// 		Label: msg,
// 		Items: []string{"Yes", "No"},
// 	}
//
// 	_, result, err := prompt.Run()
// 	if err != nil {
// 		log.Errorf("prompt failed: %s", err)
// 		return false
// 	}
//
// 	return result == "Yes"
// }
