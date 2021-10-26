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
		"logMatches": len(logMatches),
	}).Debug("filtered log matches")

	// Build a comma-separated list of discovered log file paths
	discoveredLogFiles := []string{}
	for _, logMatch := range logMatches {
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

func installAWS(ctx context.Context, i *RecipeInstaller, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error {
	log.Debugf("validating IMDS endpoint: %s", r.IMDSMatch)

	fmt.Printf("THE RECIPE STRUCT: %+v", r)

	matches := i.urlValidator.Validate(utils.SignalCtx, r.IMDSMatch)
	if !matches {
		return types.ErrorValidatingUrl
	}

	log.Debugf("validated IMDS endpoint")

	_, err := i.executeAndValidateWithProgress(ctx, m, r)
	return err
}
