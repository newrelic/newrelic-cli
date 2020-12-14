package install

type executionStatusReporter interface {
	reportUserStatus(executionStatus) error
	reportEntityStatus(string, executionStatus) error
}

type executionStatus struct {
	Timestamp   int                     `json:"timestamp"`
	Complete    bool                    `json:"complete"`
	EntityGuids []string                `json:"entityGuids"`
	Recipes     []executionStatusRecipe `json:"recipes"`
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
