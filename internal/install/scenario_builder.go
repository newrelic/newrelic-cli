package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

type TestScenario string

const (
	Basic TestScenario = "BASIC"
)

var (
	TestScenarios = []TestScenario{
		Basic,
	}
)

func TestScenarioValues() []string {
	v := make([]string, len(TestScenarios))
	for i, s := range TestScenarios {
		v[i] = string(s)
	}

	return v
}

type ScenarioBuilder struct {
	installerContext InstallerContext
}

func NewScenarioBuilder(ic InstallerContext) *ScenarioBuilder {
	b := ScenarioBuilder{
		installerContext: ic,
	}

	return &b
}

func (b *ScenarioBuilder) BuildScenario(s TestScenario) *RecipeInstaller {
	switch s {
	case Basic:
		return b.Basic()
	}

	return nil
}

func (b *ScenarioBuilder) Basic() *RecipeInstaller {

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

	i.InstallerContext = b.installerContext

	return &i
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
	f.FetchRecommendationsVal = []types.Recipe{
		{
			Name:           "Recommended recipe",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Recommended recipe
install:
  version: "3"
  tasks:
    default:
`,
		},
	}

	return f
}
