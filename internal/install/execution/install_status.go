package execution

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// nolint: maligned
type InstallStatus struct {
	Complete            bool                    `json:"complete"`
	DiscoveryManifest   types.DiscoveryManifest `json:"discoveryManifest"`
	EntityGUIDs         []string                `json:"entityGuids"`
	Error               StatusError             `json:"error"`
	LogFilePath         string                  `json:"logFilePath"`
	Statuses            []*RecipeStatus         `json:"recipes"`
	Timestamp           int64                   `json:"timestamp"`
	CLIVersion          string                  `json:"cliVersion"`
	HasInstalledRecipes bool                    `json:"hasInstalledRecipes"`
	HasCanceledRecipes  bool                    `json:"hasCanceledRecipes"`
	HasSkippedRecipes   bool                    `json:"hasSkippedRecipes"`
	HasFailedRecipes    bool                    `json:"hasFailedRecipes"`
	RecipesSkipped      []*RecipeStatus         `json:"recipesSkipped"`
	RecipesCanceled     []*RecipeStatus         `json:"recipesCanceled"`
	RecipesFailed       []*RecipeStatus         `json:"recipesFailed"`
	RecipesInstalled    []*RecipeStatus         `json:"recipesInstalled"`
	DocumentID          string
	targetedInstall     bool
	statusSubscriber    []StatusSubscriber
	successLinkConfig   types.SuccessLinkConfig
}

type RecipeStatus struct {
	DisplayName string           `json:"displayName"`
	Error       StatusError      `json:"error"`
	Name        string           `json:"name"`
	Status      RecipeStatusType `json:"status"`
	EntityGUID  string           `json:"entityGuid,omitempty"`
	// ValidationDurationMilliseconds is duration in Milliseconds that a recipe took to validate data was flowing.
	ValidationDurationMilliseconds int64 `json:"validationDurationMilliseconds,omitempty"`
}

type RecipeStatusType string

var RecipeStatusTypes = struct {
	AVAILABLE   RecipeStatusType
	CANCELED    RecipeStatusType
	INSTALLING  RecipeStatusType
	FAILED      RecipeStatusType
	INSTALLED   RecipeStatusType
	SKIPPED     RecipeStatusType
	RECOMMENDED RecipeStatusType
}{
	AVAILABLE:   "AVAILABLE",
	CANCELED:    "CANCELED",
	INSTALLING:  "INSTALLING",
	FAILED:      "FAILED",
	INSTALLED:   "INSTALLED",
	SKIPPED:     "SKIPPED",
	RECOMMENDED: "RECOMMENDED",
}

type StatusError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func NewInstallStatus(reporters []StatusSubscriber) *InstallStatus {
	s := InstallStatus{
		DocumentID:       uuid.New().String(),
		Timestamp:        utils.GetTimestamp(),
		LogFilePath:      config.DefaultConfigDirectory + "/" + config.DefaultLogFile,
		statusSubscriber: reporters,
	}

	return &s
}

func (s *InstallStatus) DiscoveryComplete(dm types.DiscoveryManifest) {
	s.withDiscoveryInfo(dm)

	for _, r := range s.statusSubscriber {
		if err := r.DiscoveryComplete(s, dm); err != nil {
			log.Errorf("Could not report discovery info: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeAvailable(recipe types.Recipe) {
	s.withAvailableRecipe(recipe)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeAvailable(s, recipe); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipesAvailable(recipes []types.Recipe) {
	s.withAvailableRecipes(recipes)

	for _, r := range s.statusSubscriber {
		if err := r.RecipesAvailable(s, recipes); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipesSelected(recipes []types.Recipe) {
	for _, r := range s.statusSubscriber {
		if err := r.RecipesSelected(s, recipes); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *InstallStatus) RecipeInstalled(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalled(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

// RecipeRecommended is responsible for setting the nerstorage scopes
// when a recipe is recommended.  This is used when a recipe is found, but not
// a "HOST" type, and is used to indicate to the user that it is something they
// should consider integrating, but not something that the recipe framework
// will currently assist with.
func (s *InstallStatus) RecipeRecommended(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.RECOMMENDED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeRecommended(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeInstalling(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.INSTALLING)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeInstalling(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeFailed(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.FAILED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeFailed(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) RecipeSkipped(event RecipeStatusEvent) {
	s.withRecipeEvent(event, RecipeStatusTypes.SKIPPED)

	for _, r := range s.statusSubscriber {
		if err := r.RecipeSkipped(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *InstallStatus) InstallComplete() {
	s.completed()

	for _, r := range s.statusSubscriber {
		if err := r.InstallComplete(s); err != nil {
			log.Errorf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) InstallCanceled() {
	s.canceled()

	for _, r := range s.statusSubscriber {
		if err := r.InstallCanceled(s); err != nil {
			log.Errorf("Error writing execution status: %s", err)
		}
	}
}

func (s *InstallStatus) recommendations() []*RecipeStatus {
	var statuses []*RecipeStatus

	for _, st := range s.Statuses {
		if st.Status == RecipeStatusTypes.RECOMMENDED {
			statuses = append(statuses, st)
		}
	}

	return statuses
}

func (s *InstallStatus) hasFailed() bool {
	for _, ss := range s.Statuses {
		if ss.Status == RecipeStatusTypes.FAILED {
			return true
		}
	}

	return false
}

func (s *InstallStatus) isCanceled() bool {
	for _, ss := range s.Statuses {
		if ss.Status == RecipeStatusTypes.CANCELED {
			return true
		}
	}

	return false
}

func (s *InstallStatus) SetTargetedInstall() {
	s.targetedInstall = true
}

func (s *InstallStatus) IsTargetedInstall() bool {
	return s.targetedInstall
}

func (s *InstallStatus) HostEntityGUID() string {
	var guid string

	// When we have performed a targeted installation, we want to roll up to the last GUID in the list.
	if len(s.EntityGUIDs) > 0 {
		if s.IsTargetedInstall() {
			guid = s.EntityGUIDs[len(s.EntityGUIDs)-1]
		} else {
			guid = s.EntityGUIDs[0]
		}
	}

	return guid
}

func (s *InstallStatus) withAvailableRecipes(recipes []types.Recipe) {
	for _, r := range recipes {
		s.withAvailableRecipe(r)
	}
}

func (s *InstallStatus) withAvailableRecipe(r types.Recipe) {
	e := RecipeStatusEvent{Recipe: r}
	s.withRecipeEvent(e, RecipeStatusTypes.AVAILABLE)
}

func (s *InstallStatus) withSuccessLinkConfig(l types.SuccessLinkConfig) {
	s.successLinkConfig = l
}

func (s *InstallStatus) withEntityGUID(entityGUID string) {
	for _, e := range s.EntityGUIDs {
		if e == entityGUID {
			return
		}
	}

	log.WithFields(log.Fields{
		"guid": entityGUID,
	}).Debug("new GUID")

	s.EntityGUIDs = append(s.EntityGUIDs, entityGUID)
}

func (s *InstallStatus) withDiscoveryInfo(dm types.DiscoveryManifest) {
	s.DiscoveryManifest = dm
	s.Timestamp = utils.GetTimestamp()

	version := os.Getenv("NEW_RELIC_CLI_VERSION")
	if version != "" {
		s.CLIVersion = version
	}
}

func (s *InstallStatus) withRecipeEvent(e RecipeStatusEvent, rs RecipeStatusType) {
	if e.EntityGUID != "" {
		s.withEntityGUID(e.EntityGUID)
	}

	s.withSuccessLinkConfig(e.Recipe.SuccessLinkConfig)

	statusError := StatusError{
		Message: e.Msg,
	}

	s.Error = statusError

	log.WithFields(log.Fields{
		"recipe_name":                    e.Recipe.Name,
		"status":                         rs,
		"error":                          statusError.Message,
		"guid":                           e.EntityGUID,
		"validationDurationMilliseconds": e.ValidationDurationMilliseconds,
	}).Debug("recipe event")

	found := s.getStatus(e.Recipe)

	if found != nil {
		found.Status = rs

		if e.EntityGUID != "" {
			found.EntityGUID = e.EntityGUID
		}

		if e.ValidationDurationMilliseconds > 0 {
			found.ValidationDurationMilliseconds = e.ValidationDurationMilliseconds
		}
	} else {
		recipeStatus := &RecipeStatus{
			Name:        e.Recipe.Name,
			DisplayName: e.Recipe.DisplayName,
			Status:      rs,
			Error:       statusError,
		}

		if e.EntityGUID != "" {
			recipeStatus.EntityGUID = e.EntityGUID
		}

		if e.ValidationDurationMilliseconds > 0 {
			recipeStatus.ValidationDurationMilliseconds = e.ValidationDurationMilliseconds
		}

		s.Statuses = append(s.Statuses, recipeStatus)
	}

	s.Timestamp = utils.GetTimestamp()
}

func (s *InstallStatus) completed() {
	s.Complete = true
	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"timestamp": s.Timestamp,
	}).Debug("completed")

	s.updateFinalInstallationStatuses(false)
}

func (s *InstallStatus) canceled() {
	s.Timestamp = utils.GetTimestamp()

	log.WithFields(log.Fields{
		"timestamp": s.Timestamp,
	}).Debug("canceled")

	s.updateFinalInstallationStatuses(true)
}

func (s *InstallStatus) getStatus(r types.Recipe) *RecipeStatus {
	for _, recipe := range s.Statuses {
		if recipe.Name == r.Name {
			return recipe
		}
	}

	return nil
}

// This function handles updating the final recipe statuses and top-level installation status.
// Canceling (e.g. ctl+c) will cause unresolved recipes to be marked as canceled.
// Exiting early (i.e. an error occurred) will cause unresolved recipes to be marked as failed.
func (s *InstallStatus) updateFinalInstallationStatuses(installCanceled bool) {
	for i, ss := range s.Statuses {
		if ss.Status == RecipeStatusTypes.AVAILABLE || ss.Status == RecipeStatusTypes.INSTALLING {
			debugMsg := "failed"

			if installCanceled {
				debugMsg = "canceled"
			}

			log.WithFields(log.Fields{
				"recipe": s.Statuses[i].Name,
			}).Debug(fmt.Sprintf("marking recipe %s", debugMsg))

			if installCanceled {
				s.Statuses[i].Status = RecipeStatusTypes.CANCELED
			} else {
				s.Statuses[i].Status = RecipeStatusTypes.FAILED
			}
		}

		// Installed
		if ss.Status == RecipeStatusTypes.INSTALLED {
			s.RecipesInstalled = append(s.RecipesInstalled, ss)
			s.HasInstalledRecipes = true
		}

		// Skipped
		if ss.Status == RecipeStatusTypes.SKIPPED {
			s.RecipesSkipped = append(s.RecipesSkipped, ss)
			s.HasSkippedRecipes = true
		}

		// Canceled
		if ss.Status == RecipeStatusTypes.CANCELED {
			s.RecipesCanceled = append(s.RecipesCanceled, ss)
			s.HasCanceledRecipes = true
		}

		// Errored
		if ss.Status == RecipeStatusTypes.FAILED {
			s.RecipesFailed = append(s.RecipesFailed, ss)
			s.HasFailedRecipes = true
		}
	}

	log.WithFields(log.Fields{
		"hasInstalledRecipes": s.HasInstalledRecipes,
		"hasSkippedRecipes":   s.HasSkippedRecipes,
		"hasCanceledRecipes":  s.HasCanceledRecipes,
		"hasFailedRecipes":    s.HasFailedRecipes,
	}).Debug("final installation statuses updated")
}
