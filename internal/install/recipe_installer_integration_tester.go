package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

type RecipeInstallerIntegrationTester struct {
	*RecipeInstaller
}

func NewRecipeInstallerIntegrationTester(ic InstallerContext) *RecipeInstallerIntegrationTester {
	t := RecipeInstallerIntegrationTester{}

	// mock implementations
	rf := setupRecipeFetcher()
	er := execution.NewMockStatusReporter()
	v := validation.NewMockRecipeValidator()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		statusReporter:    er,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = ic

	t.RecipeInstaller = &i

	return &t
}

func setupRecipeFetcher() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           "Infrastructure Agent Installer",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Infrastructure Agent Installer
install:
  version: "3"
  tasks:
    default:
`,
		},
		{
			Name:           "Logs integration",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Logs integration
install:
  version: "3"
  tasks:
    default:
`,
		},
	}

	return f
}
