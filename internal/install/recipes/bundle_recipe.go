package recipes

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	Recipe           *types.OpenInstallationRecipe
	Dependencies     []*BundleRecipe
	DetectedStatuses []*DetectedStatusType
}

type DetectedStatusType struct {
	Status     execution.RecipeStatusType
	DurationMs int64
}

func (br *BundleRecipe) AddDetectionStatus(newStatus execution.RecipeStatusType, durationMs int64) {
	if br.HasStatus(newStatus) {
		return
	}
	br.DetectedStatuses = append(br.DetectedStatuses, &DetectedStatusType{Status: newStatus, DurationMs: durationMs})
}

func (br *BundleRecipe) HasStatus(status execution.RecipeStatusType) bool {
	for _, detectedStatus := range br.DetectedStatuses {
		if detectedStatus.Status == status {
			return true
		}
	}
	return false
}

func (br *BundleRecipe) AreAllDependenciesAvailable() bool {
	for _, depName := range br.Recipe.Dependencies {
		if br.IsNameInDependencies(depName) {
			continue
		}
		log.Debugf("recipe %s is missing a dependency for %s", br.Recipe.Name, depName)
		return false
	}
	for _, ds := range br.Dependencies {
		if !ds.HasStatus(execution.RecipeStatusTypes.AVAILABLE) {
			return false
		}
	}
	return true
}

func (br *BundleRecipe) IsNameInDependencies(depName string) bool {
	for _, ds := range br.Dependencies {
		if ds.Recipe.Name == depName {
			return true
		}
	}
	return false
}

func (ds *DetectedStatusType) String() string {
	result := string(ds.Status)
	if ds.DurationMs > 0 {
		result = fmt.Sprintf("%s %dms", result, ds.DurationMs)
	}
	return fmt.Sprintf("{%s}", result)
}

func (br *BundleRecipe) String() string {
	result := fmt.Sprintf("%s %s", br.Recipe.Name, br.DetectedStatuses)
	return fmt.Sprintf("{%s}", result)
}
