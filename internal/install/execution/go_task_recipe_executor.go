package execution

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

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

func (re *GoTaskRecipeExecutor) Execute(ctx context.Context, m types.DiscoveryManifest, r types.Recipe) error {
	log.Debugf("Executing recipe %s", r.Name)

	f, err := recipes.RecipeToRecipeFile(r)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(f.Install)
	if err != nil {
		return err
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
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	if err = e.Setup(); err != nil {
		return err
	}

	var tf taskfile.Taskfile
	err = yaml.Unmarshal(out, &tf)
	if err != nil {
		return err
	}

	calls, globals := taskargs.ParseV3()
	e.Taskfile.Vars.Merge(globals)

	setVarsFromSystemInfo(e.Taskfile, m)

	if err := setVarsFromProfile(e.Taskfile); err != nil {
		return err
	}

	if err := setVarsFromInput(e.Taskfile, f.InputVars); err != nil {
		return err
	}

	if err := setStaticVars(e.Taskfile, r.Vars); err != nil {
		return err
	}

	if err := e.Run(ctx, calls...); err != nil {
		return err
	}

	return nil
}

func setVarsFromProfile(t *taskfile.Taskfile) error {
	defaultProfile := credentials.DefaultProfile()
	if defaultProfile.LicenseKey == "" {
		return errors.New("license key not found in default profile")
	}

	v := taskfile.Vars{}
	v.Set("NEW_RELIC_LICENSE_KEY", taskfile.Var{Static: defaultProfile.LicenseKey})
	v.Set("NEW_RELIC_ACCOUNT_ID", taskfile.Var{Static: strconv.Itoa(defaultProfile.AccountID)})
	v.Set("NEW_RELIC_API_KEY", taskfile.Var{Static: defaultProfile.APIKey})
	v.Set("NEW_RELIC_REGION", taskfile.Var{Static: defaultProfile.Region})

	t.Vars.Merge(&v)

	return nil
}

func setVarsFromSystemInfo(t *taskfile.Taskfile, m types.DiscoveryManifest) {
	v := taskfile.Vars{}
	v.Set("HOSTNAME", taskfile.Var{Static: m.Hostname})
	v.Set("OS", taskfile.Var{Static: m.OS})
	v.Set("PLATFORM", taskfile.Var{Static: m.Platform})
	v.Set("PLATFORM_FAMILY", taskfile.Var{Static: m.PlatformFamily})
	v.Set("PLATFORM_VERSION", taskfile.Var{Static: m.PlatformVersion})
	v.Set("KERNEL_ARCH", taskfile.Var{Static: m.KernelArch})
	v.Set("KERNEL_VERSION", taskfile.Var{Static: m.KernelVersion})

	t.Vars.Merge(&v)
}

func setStaticVars(t *taskfile.Taskfile, vars map[string]interface{}) error {
	for k, x := range vars {

		varData, err := yaml.Marshal(x)
		if err != nil {
			return err
		}

		t.Vars.Set(k, taskfile.Var{Static: string(varData)})
	}

	return nil
}

func setVarsFromInput(t *taskfile.Taskfile, inputVars []recipes.VariableConfig) error {
	for _, envConfig := range inputVars {
		v := taskfile.Vars{}

		envValue := os.Getenv(envConfig.Name)
		if envValue == "" {
			log.Debugf("required env var %s not found", envConfig.Name)
			msg := fmt.Sprintf("value for %s required", envConfig.Name)

			if envConfig.Prompt != "" {
				msg = envConfig.Prompt
			}

			prompt := promptui.Prompt{
				Label: msg,
			}

			if envConfig.Default != "" {
				prompt.Default = envConfig.Default
			}

			result, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %s", err)
			}

			v.Set(envConfig.Name, taskfile.Var{Static: result})
		} else {
			v.Set(envConfig.Name, taskfile.Var{Static: envValue})
		}

		t.Vars.Merge(&v)
	}

	return nil
}
