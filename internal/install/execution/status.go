package execution

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type StatusRollup struct {
	Complete        bool `json:"complete"`
	DocumentID      string
	EntityGUIDs     []string `json:"entityGuids"`
	Statuses        []Status `json:"recipes"`
	Timestamp       int64    `json:"timestamp"`
	LogFilePath     string   `json:"logFilePath"`
	statusReporters []StatusReporter
}

type Status struct {
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	Status      StatusType          `json:"status"`
	Errors      []StatusRecipeError `json:"errors"`
}

type StatusType string

var StatusTypes = struct {
	AVAILABLE  StatusType
	INSTALLING StatusType
	FAILED     StatusType
	INSTALLED  StatusType
	SKIPPED    StatusType
}{
	AVAILABLE:  "AVAILABLE",
	INSTALLING: "INSTALLING",
	FAILED:     "FAILED",
	INSTALLED:  "INSTALLED",
	SKIPPED:    "SKIPPED",
}

type StatusRecipeError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func NewStatusRollup(reporters []StatusReporter) *StatusRollup {

	s := StatusRollup{
		DocumentID:      uuid.New().String(),
		Timestamp:       utils.GetTimestamp(),
		LogFilePath:     config.DefaultConfigDirectory + "/" + config.DefaultLogFile,
		statusReporters: reporters,
	}

	return &s
}

func (s *StatusRollup) ReportRecipeAvailable(recipe types.Recipe) {
	s.withAvailableRecipe(recipe)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipeAvailable(s, recipe); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *StatusRollup) ReportRecipesAvailable(recipes []types.Recipe) {
	s.withAvailableRecipes(recipes)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipesAvailable(s, recipes); err != nil {
			log.Errorf("Could not report recipe execution status: %s", err)
		}
	}
}

func (s *StatusRollup) ReportRecipeInstalled(event RecipeStatusEvent) {
	s.withRecipeEvent(event, StatusTypes.INSTALLED)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipeInstalled(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *StatusRollup) ReportRecipeInstalling(event RecipeStatusEvent) {
	s.withRecipeEvent(event, StatusTypes.INSTALLING)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipeInstalling(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *StatusRollup) ReportRecipeFailed(event RecipeStatusEvent) {
	s.withRecipeEvent(event, StatusTypes.FAILED)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipeFailed(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *StatusRollup) ReportRecipeSkipped(event RecipeStatusEvent) {
	s.withRecipeEvent(event, StatusTypes.SKIPPED)

	for _, r := range s.statusReporters {
		if err := r.ReportRecipeSkipped(s, event); err != nil {
			log.Errorf("Error writing recipe status for recipe %s: %s", event.Recipe.Name, err)
		}
	}
}

func (s *StatusRollup) ReportComplete() {
	s.Complete = true
	s.Timestamp = utils.GetTimestamp()

	for _, r := range s.statusReporters {
		if err := r.ReportComplete(s); err != nil {
			log.Errorf("Error writing execution status: %s", err)
		}
	}
}

func (s *StatusRollup) hasFailed() bool {
	for _, ss := range s.Statuses {
		if ss.Status == StatusTypes.FAILED {
			return true
		}
	}

	return false
}

func (s *StatusRollup) withAvailableRecipes(recipes []types.Recipe) {
	for _, r := range recipes {
		s.withAvailableRecipe(r)
	}
}

func (s *StatusRollup) withAvailableRecipe(r types.Recipe) {
	e := RecipeStatusEvent{Recipe: r}
	s.withRecipeEvent(e, StatusTypes.AVAILABLE)
}

func (s *StatusRollup) withEntityGUID(entityGUID string) {
	for _, e := range s.EntityGUIDs {
		if e == entityGUID {
			return
		}
	}

	s.EntityGUIDs = append(s.EntityGUIDs, entityGUID)
}

func (s *StatusRollup) withRecipeEvent(e RecipeStatusEvent, rs StatusType) {
	if e.EntityGUID != "" {
		s.withEntityGUID(e.EntityGUID)
	}

	found := s.getStatus(e.Recipe)

	if found != nil {
		found.Status = rs
	} else {
		e := &Status{
			Name:        e.Recipe.Name,
			DisplayName: e.Recipe.DisplayName,
			Status:      rs,
		}
		s.Statuses = append(s.Statuses, *e)
	}

	s.Timestamp = utils.GetTimestamp()
}

func (s *StatusRollup) getStatus(r types.Recipe) *Status {
	var found *Status
	for i, recipe := range s.Statuses {
		if recipe.Name == r.Name {
			found = &s.Statuses[i]
		}
	}

	return found
}
