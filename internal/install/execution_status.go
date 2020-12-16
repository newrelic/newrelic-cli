package install

import (
	"time"

	"github.com/google/uuid"
)

type executionStatus struct {
	Complete    bool `json:"complete"`
	DocumentID  string
	EntityGuids []string                `json:"entityGuids"`
	Recipes     []executionStatusRecipe `json:"recipes"`
	Timestamp   int64                   `json:"timestamp"`
}

type executionStatusRecipe struct {
	Name        string                       `json:"name"`
	DisplayName string                       `json:"displayName"`
	Status      executionStatusRecipeStatus  `json:"status"`
	Errors      []executionStatusRecipeError `json:"errors"`
}

type executionStatusRecipeStatus string

var executionStatusRecipeStatusTypes = struct {
	AVAILABLE executionStatusRecipeStatus
	FAILED    executionStatusRecipeStatus
	INSTALLED executionStatusRecipeStatus
	SKIPPED   executionStatusRecipeStatus
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

func newExecutionStatus() executionStatus {
	s := executionStatus{
		DocumentID: uuid.New().String(),
		Timestamp:  getTimestamp(),
	}

	return s
}

func getTimestamp() int64 {
	return time.Now().Unix()
}

func (s *executionStatus) withAvailableRecipes(recipes []recipe) {
	for _, r := range recipes {
		e := recipeStatusEvent{recipe: r}
		s.withRecipeEvent(e, executionStatusRecipeStatusTypes.AVAILABLE)
	}
}

func (s *executionStatus) withRecipeEvent(e recipeStatusEvent, rs executionStatusRecipeStatus) {
	found := s.getExecutionStatusRecipe(e.recipe)

	if found != nil {
		found.Status = rs
	} else {
		e := &executionStatusRecipe{
			Name:        e.recipe.Name,
			DisplayName: e.recipe.DisplayName,
			Status:      rs,
		}
		s.Recipes = append(s.Recipes, *e)
	}

	s.Timestamp = getTimestamp()
}

func (s *executionStatus) getExecutionStatusRecipe(r recipe) *executionStatusRecipe {
	var found *executionStatusRecipe
	for i, recipe := range s.Recipes {
		if recipe.Name == r.Name {
			found = &s.Recipes[i]
		}
	}

	return found
}
