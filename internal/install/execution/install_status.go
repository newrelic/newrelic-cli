package execution

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type InstallStatus struct {
	Complete          bool                    `json:"complete"`
	DiscoveryManifest types.DiscoveryManifest `json:"discoveryManifest"`
	EntityGUIDs       []string                `json:"entityGuids"`
	Error             StatusError             `json:"error"`
	LogFilePath       string                  `json:"logFilePath"`
	Statuses          []RecipeStatus          `json:"recipes"`
	Timestamp         int64                   `json:"timestamp"`
	DocumentID        string
	statusSubscriber  []StatusSubscriber
}

type RecipeStatus struct {
	DisplayName string           `json:"displayName"`
	Error       StatusError      `json:"error"`
	Name        string           `json:"name"`
	Status      RecipeStatusType `json:"status"`
}

type RecipeStatusType string

var RecipeStatusTypes = struct {
	AVAILABLE   RecipeStatusType
	INSTALLING  RecipeStatusType
	FAILED      RecipeStatusType
	INSTALLED   RecipeStatusType
	SKIPPED     RecipeStatusType
	RECOMMENDED RecipeStatusType
}{
	AVAILABLE:   "AVAILABLE",
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

func (s *InstallStatus) recommendations() []RecipeStatus {
	var statuses []RecipeStatus

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

func (s *InstallStatus) withAvailableRecipes(recipes []types.Recipe) {
	for _, r := range recipes {
		s.withAvailableRecipe(r)
	}
}

func (s *InstallStatus) withAvailableRecipe(r types.Recipe) {
	e := RecipeStatusEvent{Recipe: r}
	s.withRecipeEvent(e, RecipeStatusTypes.AVAILABLE)
}

func (s *InstallStatus) withEntityGUID(entityGUID string) {
	for _, e := range s.EntityGUIDs {
		if e == entityGUID {
			return
		}
	}

	s.EntityGUIDs = append(s.EntityGUIDs, entityGUID)
}

func (s *InstallStatus) withDiscoveryInfo(dm types.DiscoveryManifest) {
	s.DiscoveryManifest = dm
	s.Timestamp = utils.GetTimestamp()
}

func (s *InstallStatus) withRecipeEvent(e RecipeStatusEvent, rs RecipeStatusType) {
	if e.EntityGUID != "" {
		s.withEntityGUID(e.EntityGUID)
	}

	statusError := StatusError{
		Message: e.Msg,
	}

	s.Error = statusError

	log.WithFields(log.Fields{
		"recipe_name": e.Recipe.Name,
		"status":      rs,
		"error":       statusError.Message,
	}).Debug("recipe event")

	found := s.getStatus(e.Recipe)

	if found != nil {
		found.Status = rs
	} else {
		e := &RecipeStatus{
			Name:        e.Recipe.Name,
			DisplayName: e.Recipe.DisplayName,
			Status:      rs,
			Error:       statusError,
		}
		s.Statuses = append(s.Statuses, *e)
	}

	s.Timestamp = utils.GetTimestamp()
}

func (s *InstallStatus) completed() {
	s.Complete = true
	s.Timestamp = utils.GetTimestamp()

	// Exiting early will cause unresolved recipes to be marked as failed.
	for i, ss := range s.Statuses {
		if ss.Status == RecipeStatusTypes.AVAILABLE || ss.Status == RecipeStatusTypes.INSTALLING {
			s.Statuses[i].Status = RecipeStatusTypes.FAILED
		}
	}
}

func (s *InstallStatus) getStatus(r types.Recipe) *RecipeStatus {
	var found *RecipeStatus
	for i, recipe := range s.Statuses {
		if recipe.Name == r.Name {
			found = &s.Statuses[i]
		}
	}

	return found
}
