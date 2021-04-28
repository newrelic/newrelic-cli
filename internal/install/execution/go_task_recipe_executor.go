package execution

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// GoTaskRecipeExecutor is an implementation of the recipeExecutor interface that
// uses the go-task module to execute the steps defined in each recipe.
type GoTaskRecipeExecutor struct{}

// NewGoTaskRecipeExecutor returns a new instance of GoTaskRecipeExecutor.
func NewGoTaskRecipeExecutor() *GoTaskRecipeExecutor {
	return &GoTaskRecipeExecutor{}
}

func (re *GoTaskRecipeExecutor) Prepare(ctx context.Context, m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error) {
	log.WithFields(log.Fields{
		"name": r.Name,
	}).Debug("preparing recipe")

	vars := types.RecipeVars{}

	results := []types.RecipeVars{}

	systemInfoResult := varsFromSystemInfo(m)

	profileResult, err := varsFromProfile(licenseKey)
	if err != nil {
		return types.RecipeVars{}, err
	}

	inputVarsResult, err := varsFromInput(r.InputVars, assumeYes)
	if err != nil {
		return types.RecipeVars{}, err
	}

	results = append(results, systemInfoResult)
	results = append(results, profileResult)
	results = append(results, types.RecipeVariables)
	results = append(results, inputVarsResult)

	for _, result := range results {
		for k, v := range result {
			vars[k] = v
		}
	}

	return vars, nil
}

func (re *GoTaskRecipeExecutor) Execute(ctx context.Context, m types.DiscoveryManifest, r types.OpenInstallationRecipe, recipeVars types.RecipeVars) error {
	log.Debugf("executing recipe %s", r.Name)

	out := []byte(r.Install)

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
		Stdin:      os.Stdin,
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
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("Task execution returned error")

		// go-task does not provide an error type to denote context cancelation
		// Therefore we need to match inside the error message
		if strings.Contains(err.Error(), "context canceled") {
			return types.ErrInterrupt
		}

		// Recipe trigger a canceled event with specific exit code 130 used
		if strings.Contains(err.Error(), "exit status 130") {
			return types.ErrInterrupt
		}

		return err
	}

	return nil
}

func varsFromProfile(licenseKey string) (types.RecipeVars, error) {
	defaultProfile := credentials.DefaultProfile()
	if licenseKey == "" {
		return types.RecipeVars{}, errors.New("license key not found")
	}

	vars := make(types.RecipeVars)

	vars["NEW_RELIC_LICENSE_KEY"] = licenseKey
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

func varsFromInput(inputVars []types.OpenInstallationRecipeInputVariable, assumeYes bool) (types.RecipeVars, error) {
	vars := make(types.RecipeVars)

	vars["NEW_RELIC_ASSUME_YES"] = fmt.Sprintf("%t", assumeYes)

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
				if err == terminal.InterruptErr {
					return types.RecipeVars{}, types.ErrInterrupt
				}

				return types.RecipeVars{}, fmt.Errorf("prompt failed: %s", err)
			}
		}

		vars[envConfig.Name] = envValue
	}

	return vars, nil
}

func varFromPrompt(envConfig types.OpenInstallationRecipeInputVariable) (string, error) {
	msg := fmt.Sprintf("value for %s required", envConfig.Name)

	if envConfig.Prompt != "" {
		msg = envConfig.Prompt
	}

	value := ""
	var prompt survey.Prompt

	if envConfig.Secret {
		prompt = &survey.Password{
			Message: msg,
		}
	} else {
		defaultValue := ""

		if envConfig.Default != "" {
			defaultValue = envConfig.Default
		}

		prompt = &survey.Input{
			Message: msg,
			Default: defaultValue,
		}
	}

	err := survey.AskOne(prompt, &value)
	if err != nil {
		return "", err
	}

	return value, nil

}
