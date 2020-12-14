package install

type executionStatusReporter interface {
	reportUserStatus(executionStatus) error
	reportEntityStatus(executionStatus) error
}

type executionStatus struct {
	Timestamp   int                     `json:"timestamp"`
	Complete    bool                    `json:"complete"`
	EntityGuids []string                `json:"entityGuids"`
	Recipes     []executionStatusRecipe `json:"recipes"`
}

type executionStatusRecipe struct {
	Name        string                      `json:"name"`
	DisplayName string                      `json:"displayName"`
	Status      executionStatusRecipeStatus `json:"status"`
	Errors      []string                    `json:"errors"`
}

type executionStatusRecipeStatus string

var executionStatusRecipeStatusTypes = struct {
	AVAILABLE executionStatusRecipeStatus
	INSTALLED executionStatusRecipeStatus
	FAILED    executionStatusRecipeStatus
	SKIPPED   executionStatusRecipeStatus
}{
	AVAILABLE: "AVAILABLE",
	INSTALLED: "INSTALLED",
	SKIPPED:   "SKIPPED",
}
