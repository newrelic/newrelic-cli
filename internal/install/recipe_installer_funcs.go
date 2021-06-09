package install

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func installLogging(ctx context.Context, i *RecipeInstaller, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error {
	log.WithFields(log.Fields{
		"recipe_count": len(recipes),
	}).Debug("filtering log matches")
	logMatches, err := i.fileFilterer.Filter(utils.SignalCtx, recipes)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"possible_matches": len(logMatches),
	}).Debug("filtered log matches")

	var acceptedLogMatches []types.OpenInstallationLogMatch
	var ok bool
	for _, match := range logMatches {
		ok, err = i.userAcceptsLogFile(match)
		if err != nil {
			return err
		}

		if ok {
			acceptedLogMatches = append(acceptedLogMatches, match)
		}
	}

	log.WithFields(log.Fields{
		"matches": acceptedLogMatches,
	}).Debug("matches accepted")

	// Build a comma-separated list of discovered log file paths
	discoveredLogFiles := []string{}
	for _, logMatch := range acceptedLogMatches {
		discoveredLogFiles = append(discoveredLogFiles, logMatch.File)
	}

	discoveredLogFilesString := strings.Join(discoveredLogFiles, ",")
	r.SetRecipeVar("NR_DISCOVERED_LOG_FILES", discoveredLogFilesString)

	log.WithFields(log.Fields{
		"NR_DISCOVERED_LOG_FILES": discoveredLogFilesString,
	}).Debug("discovered log files")

	_, err = i.executeAndValidateWithProgress(ctx, m, r)
	return err
}

func (i *RecipeInstaller) userAccepts(msg string) (bool, error) {
	if i.AssumeYes {
		return true, nil
	}

	val, err := i.prompter.PromptYesNo(msg)
	if err != nil {
		return false, err
	}

	return val, nil
}

func (i *RecipeInstaller) userAcceptsLogFile(match types.OpenInstallationLogMatch) (bool, error) {
	msg := fmt.Sprintf("Files have been found at the following pattern: %s Do you want to watch them?", match.File)
	return i.userAccepts(msg)
}
