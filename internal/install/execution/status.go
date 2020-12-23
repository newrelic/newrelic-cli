package execution

import (
	"github.com/google/uuid"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type StatusRollup struct {
	Complete    bool `json:"complete"`
	DocumentID  string
	EntityGUIDs []string `json:"entityGuids"`
	Statuses    []Status `json:"recipes"`
	Timestamp   int64    `json:"timestamp"`
}

type Status struct {
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	Status      StatusType          `json:"status"`
	Errors      []StatusRecipeError `json:"errors"`
}

type StatusType string

var StatusTypes = struct {
	AVAILABLE StatusType
	FAILED    StatusType
	INSTALLED StatusType
	SKIPPED   StatusType
}{
	AVAILABLE: "AVAILABLE",
	FAILED:    "FAILED",
	INSTALLED: "INSTALLED",
	SKIPPED:   "SKIPPED",
}

type StatusRecipeError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func NewStatusRollup() StatusRollup {
	s := StatusRollup{
		DocumentID: uuid.New().String(),
		Timestamp:  utils.GetTimestamp(),
	}

	return s
}

func (s *StatusRollup) withAvailableRecipes(recipes []types.Recipe) {
	for _, r := range recipes {
		e := RecipeStatusEvent{Recipe: r}
		s.withRecipeEvent(e, StatusTypes.AVAILABLE)
	}
}

func (s *StatusRollup) withEntityGuid(entityGUID string) {
	for _, e := range s.EntityGUIDs {
		if e == entityGUID {
			return
		}
	}

	s.EntityGUIDs = append(s.EntityGUIDs, entityGUID)
}

func (s *StatusRollup) withRecipeEvent(e RecipeStatusEvent, rs StatusType) {
	if e.EntityGUID != "" {
		s.withEntityGuid(e.EntityGUID)
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
