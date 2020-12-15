package install

import "time"

type executionStatusReporter interface {
	reportUserStatus(executionStatus) error
	reportEntityStatus(string, executionStatus) error
}

type executionStatus struct {
	Complete    bool `json:"complete"`
	DocumentID  string
	EntityGuids []string                `json:"entityGuids"`
	Recipes     []executionStatusRecipe `json:"recipes"`
	Scope       string
	Timestamp   int64 `json:"timestamp"`
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

func newExecutionStatus(recipes []recipe) executionStatus {
	s := executionStatus{
		Timestamp: time.Now().Unix(),
		Complete:  false,
	}

	availableRecipes := []executionStatusRecipe{}

	for _, r := range recipes {
		availableRecipe := executionStatusRecipe{
			// TODO work out the details about dispay name vs short name
			Name:        r.Name,
			DisplayName: r.Name,
			Status:      executionStatusRecipeStatusTypes.AVAILABLE,
		}

		availableRecipes = append(availableRecipes, availableRecipe)
	}

	return s
}
