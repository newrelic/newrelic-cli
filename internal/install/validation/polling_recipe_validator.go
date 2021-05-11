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
func (m *PollingRecipeValidator) ValidateRecipe(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe) (string, error) {
	query, err := substituteHostname(dm, r)
	if err != nil {
		return "", err
	}

	return m.Validate(ctx, query)
}

func substituteHostname(dm types.DiscoveryManifest, r types.OpenInstallationRecipe) (string, error) {
	tmpl, err := template.New("validationNRQL").Parse(string(r.ValidationNRQL))
	if err != nil {
		panic(err)
	}

	v := struct {
		HOSTNAME string
	}{
		HOSTNAME: dm.Hostname,
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, v); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
