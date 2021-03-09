package execution

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// GoTaskRecipeExecutor is an implementation of the recipeExecutor interface that
// uses the go-task module to execute the steps defined in each recipe.
type GoTaskRecipeExecutor struct{}

// NewGoTaskRecipeExecutor returns a new instance of GoTaskRecipeExecutor.
func NewGoTaskRecipeExecutor() *GoTaskRecipeExecutor {
	return &GoTaskRecipeExecutor{}
}

func (re *GoTaskRecipeExecutor) Prepare(ctx context.Context, m types.DiscoveryManifest, r types.Recipe, assumeYes bool) (types.RecipeVars, error) {
	log.WithFields(log.Fields{
		"name": r.Name,
	}).Debug("preparing recipe")

	vars := types.RecipeVars{}

	results := []types.RecipeVars{}

	systemInfoResult := varsFromSystemInfo(m)

	profileResult, err := varsFromProfile()
	if err != nil {
		return types.RecipeVars{}, err
	}

	recipeResult, err := varsFromRecipe(r)
	if err != nil {
		return types.RecipeVars{}, err
	}

	f, err := recipes.RecipeToRecipeFile(r)
	if err != nil {
		return types.RecipeVars{}, err
	}

	inputVarsResult, err := varsFromInput(f.InputVars, assumeYes)
	if err != nil {
		return types.RecipeVars{}, err
	}

	results = append(results, systemInfoResult)
	results = append(results, profileResult)
	results = append(results, recipeResult)
	results = append(results, inputVarsResult)

	for _, result := range results {
		for k, v := range result {
			vars[k] = v
		}
	}

	return vars, nil
}

func (re *GoTaskRecipeExecutor) Execute(ctx context.Context, m types.DiscoveryManifest, r types.Recipe, recipeVars types.RecipeVars) error {
	log.Debugf("executing recipe %s", r.Name)

	f, err := recipes.RecipeToRecipeFile(r)
	if err != nil {
		return fmt.Errorf("could not convert recipe to recipe file: %s", err)
	}

	out, err := yaml.Marshal(f.Install)
	if err != nil {
		return fmt.Errorf("could not marshal recipe file: %s", err)
	}

	// Create a temporary task file.
	file, err := ioutil.TempFile("", r.Name)
	defer os.Remove(file.Name())
	if err != nil {
		return err
	}

	_, err = file.Write(out)
	if err != nil {
		return err
	}

	e := task.Executor{
		Entrypoint: file.Name(),
		Stderr:     os.Stderr,
		Stdout:     os.Stdout,
	}

	if err = e.Setup(); err != nil {
		return fmt.Errorf("could not set up task executor: %s", err)
	}

	var tf taskfile.Taskfile
	err = yaml.Unmarshal(out, &tf)
	if err != nil {
		return fmt.Errorf("could not unmarshal taskfile: %s", err)
	}

	calls, globals := taskargs.ParseV3()
	e.Taskfile.Vars.Merge(globals)

	for k, val := range recipeVars {
		e.Taskfile.Vars.Set(k, taskfile.Var{Static: val})
	}

	if err := e.Run(ctx, calls...); err != nil {
		// go-task does not provide an error type to denote context cancelation
		// Therefore we need to match inside the error message
		if strings.Contains(err.Error(), "context canceled") {
			return types.ErrInterrupt
		}

		return err
	}

	return nil
}

func varsFromProfile() (types.RecipeVars, error) {
	defaultProfile := credentials.DefaultProfile()
	if defaultProfile.LicenseKey == "" {
		return types.RecipeVars{}, errors.New("license key not found in default profile")
	}

	vars := make(types.RecipeVars)

	vars["NEW_RELIC_LICENSE_KEY"] = defaultProfile.LicenseKey
	vars["NEW_RELIC_ACCOUNT_ID"] = strconv.Itoa(defaultProfile.AccountID)
	vars["NEW_RELIC_API_KEY"] = defaultProfile.APIKey
	vars["NEW_RELIC_REGION"] = defaultProfile.Region

	return vars, nil
}

func varsFromSystemInfo(m types.DiscoveryManifest) types.RecipeVars {
	vars := make(types.RecipeVars)

	vars["HOSTNAME"] = m.Hostname
	vars["OS"] = m.OS
	vars["PLATFORM"] = m.Platform
	vars["PLATFORM_FAMILY"] = m.PlatformFamily
	vars["PLATFORM_VERSION"] = m.PlatformVersion
	vars["KERNEL_ARCH"] = m.KernelArch
	vars["KERNEL_VERSION"] = m.KernelVersion

	return vars
}

func varsFromRecipe(r types.Recipe) (types.RecipeVars, error) {
	vars := make(types.RecipeVars)

	for k, x := range r.Vars {
		varData, err := yaml.Marshal(x)
		if err != nil {
			return types.RecipeVars{}, err
		}

		vars[k] = string(varData)
	}

	return vars, nil
}

func varsFromInput(inputVars []recipes.VariableConfig, assumeYes bool) (types.RecipeVars, error) {
	vars := make(types.RecipeVars)

	for _, envConfig := range inputVars {
		var err error
		envValue := os.Getenv(envConfig.Name)

		if envValue != "" {
			vars[envConfig.Name] = envValue
			continue
		}

		if assumeYes {
			if envConfig.Default == "" {
				return types.RecipeVars{}, fmt.Errorf("no default value for environment variable %s and none provided", envConfig.Name)
			}

			log.WithFields(log.Fields{
				"name":    envConfig.Name,
				"default": envConfig.Default,
			}).Debug("required env var not found, using default")

			envValue = envConfig.Default
		} else {
			log.WithFields(log.Fields{
				"name": envConfig.Name,
			}).Debug("required environment variable not found")

			envValue, err = varFromPrompt(envConfig)
			if err != nil {
				if err == promptui.ErrInterrupt {
					return types.RecipeVars{}, types.ErrInterrupt
				}

				return types.RecipeVars{}, fmt.Errorf("prompt failed: %s", err)
			}
		}

		vars[envConfig.Name] = envValue
	}

	return vars, nil
}

func varFromPrompt(envConfig recipes.VariableConfig) (string, error) {
	msg := fmt.Sprintf("value for %s required", envConfig.Name)

	if envConfig.Prompt != "" {
		msg = envConfig.Prompt
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }} ",
		Valid:   "{{ . | bold }} ",
		Invalid: "{{ . | bold }} ",
		Success: "  - {{ . }} ",
	}

	prompt := promptui.Prompt{
		Label:     msg,
		Templates: templates,
	}

	if envConfig.Secret {
		prompt.HideEntered = true
		prompt.Mask = '*'
	}

	if envConfig.Default != "" {
		prompt.Default = envConfig.Default
	}

	return prompt.Run()
}
