package install

import (
	"time"

	"github.com/google/uuid"
)

type executionStatusRollup struct {
	Complete    bool `json:"complete"`
	DocumentID  string
	EntityGuids []string          `json:"entityGuids"`
	Statuses    []executionStatus `json:"recipes"`
	Timestamp   int64             `json:"timestamp"`
}

type executionStatus struct {
	Name        string                       `json:"name"`
	DisplayName string                       `json:"displayName"`
	Status      executionStatusType          `json:"status"`
	Errors      []executionStatusRecipeError `json:"errors"`
}

type executionStatusType string

var executionStatusTypes = struct {
	AVAILABLE executionStatusType
	FAILED    executionStatusType
	INSTALLED executionStatusType
	SKIPPED   executionStatusType
}{
	AVAILABLE: "AVAILABLE",
	FAILED:    "FAILED",
	INSTALLED: "INSTALLED",
	SKIPPED:   "SKIPPED",
}

type executionStatusRecipeError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func newExecutionStatusRollup() executionStatusRollup {
	s := executionStatusRollup{
		DocumentID: uuid.New().String(),
		Timestamp:  getTimestamp(),
	}

	return s
}

func getTimestamp() int64 {
	return time.Now().Unix()
}

func (s *executionStatusRollup) withAvailableRecipes(recipes []recipe) {
	for _, r := range recipes {
		e := recipeStatusEvent{recipe: r}
		s.withRecipeEvent(e, executionStatusTypes.AVAILABLE)
	}
}

func (s *executionStatusRollup) withRecipeEvent(e recipeStatusEvent, rs executionStatusType) {
	found := s.getExecutionStatusRecipe(e.recipe)

	if found != nil {
		found.Status = rs
	} else {
		e := &executionStatus{
			Name:        e.recipe.Name,
			DisplayName: e.recipe.DisplayName,
			Status:      rs,
		}
		s.Statuses = append(s.Statuses, *e)
	}

	s.Timestamp = getTimestamp()
}

func (s *executionStatusRollup) getExecutionStatusRecipe(r recipe) *executionStatus {
	var found *executionStatus
	for i, recipe := range s.Statuses {
		if recipe.Name == r.Name {
			found = &s.Statuses[i]
		}
	}

	return found
}
