package validation

import (
	"bytes"
	"context"
	"html/template"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	utilsValidation "github.com/newrelic/newrelic-cli/internal/utils/validation"
)

type contextKey int

// PollingRecipeValidator is an implementation of the RecipeValidator interface
// that polls NRDB to assert data is being reported for the given recipe.
type PollingRecipeValidator struct {
	utilsValidation.PollingNRQLValidator
}

// NewPollingRecipeValidator returns a new instance of PollingRecipeValidator.
func NewPollingRecipeValidator(c utils.NRDBClient) *PollingRecipeValidator {
	v := PollingRecipeValidator{
		PollingNRQLValidator: *utilsValidation.NewPollingNRQLValidator(c),
	}

	return &v
}

// ValidateRecipe polls NRDB to assert data is being reported for the given recipe.
func (m *PollingRecipeValidator) ValidateRecipe(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	query, err := substituteRecipeVars(dm, r, vars)
	if err != nil {
		return "", err
	}

	return m.Validate(ctx, query)
}

func substituteRecipeVars(dm types.DiscoveryManifest, r types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	tmpl, err := template.New("validationNRQL").Parse(string(r.ValidationNRQL))
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, vars); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
