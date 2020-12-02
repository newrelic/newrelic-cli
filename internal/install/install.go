package install

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

func install(client *newrelic.NewRelic) error {
	rf := newServiceRecipeFetcher(&client.NerdGraph)
	pf := newRegexProcessFilterer(rf)

	// Execute the discovery process.
	log.Debug("Running discovery...")
	var d discoverer = newPSUtilDiscoverer(pf)
	m, err := d.discover(utils.SignalCtx)
	if err != nil {
		return err
	}

	// Retrieve the relevant recipes.
	log.Debug("Retrieving recipes...")
	recipes, err := rf.fetchRecommendations(utils.SignalCtx, m)
	if err != nil {
		return err
	}

	// Use the received recipies to determine the log file locations, prompt the
	// user for acceptance.
	var logMatches []logMatcher
	var logFiles []string
	for _, r := range recipes {
		for _, l := range r.MELTMatch.Logging {
			match, files := matchLogFilesFromRecipe(l)
			if match {
				if userAcceptLogFiles(files) {
					logMatches = append(logMatches, l)
				}
			}
		}
	}

	// LOG_FILES is the name of the variable used by the logging recipe as input
	// for creating the file.  When empty, we use the results of the recipes to
	// determine the files.
	if os.Getenv("LOG_FILES") != "" && len(logFiles) > 0 {
		os.Setenv("LOG_FILES", strings.Join(logFiles, ","))
	}

	// Execute the recipe steps.
	for _, r := range recipes {
		err := executeRecipeSteps(utils.SignalCtx, *m, r)
		if err != nil {
			return err
		}
	}

	return nil
}

// matchLogFilesFromRecipe determines if any files match the given logMatcher.
func matchLogFilesFromRecipe(matcher logMatcher) (bool, []string) {
	matches, err := filepath.Glob(matcher.File)
	if err != nil {
		log.Errorf("error matching logfiles: %s", err)
		return false, nil
	}

	if len(matches) > 0 {
		return true, matches
	}

	return false, nil
}

func userAcceptLogFiles(files []string) bool {
	msg := fmt.Sprintf("The following log files have been found: %s\nDo you want to watch them? [Yes/No]", strings.Join(files, ", "))

	prompt := promptui.Select{
		Label: msg,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Errorf("prompt failed: %s", err)
		return false
	}

	return result == "Yes"
}

func matchLog(logMatch logMatcher) bool {
	matches, err := filepath.Glob(logMatch.File)
	if err != nil {
		log.Errorf("error matching logfiles: %s", err)
		return false
	}

	if len(matches) > 0 {
		userAcceptLogFiles(matches)
	}

	return false
}

func executeRecipeSteps(ctx context.Context, m discoveryManifest, r recipeFile) error {
	log.Debugf("Executing recipe %s", r.Name)

	out, err := yaml.Marshal(r.Install)
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

	credentials.WithCredentials(func(c *credentials.Credentials) {
		v := taskfile.Vars{}
		licenseKey := c.Profiles[c.DefaultProfile].LicenseKey
		if licenseKey == "" {
			err = errors.New("license key not found in default profile")
		}

		v.Set("NR_LICENSE_KEY", taskfile.Var{Static: licenseKey})
		e.Taskfile.Vars.Merge(&v)
	})

	if err != nil {
		return err
	}

	setSystemVars(e.Taskfile, m)

	if err := setInputVars(e.Taskfile, r.InputVars); err != nil {
		return err
	}

	if err := e.Run(ctx, calls...); err != nil {
		return err
	}

	return nil
}

func setSystemVars(t *taskfile.Taskfile, m discoveryManifest) {
	v := taskfile.Vars{}
	v.Set("OS", taskfile.Var{Static: m.os})
	v.Set("Platform", taskfile.Var{Static: m.platform})
	v.Set("PlatformFamily", taskfile.Var{Static: m.platformFamily})
	v.Set("PlatformVersion", taskfile.Var{Static: m.platformVersion})
	v.Set("KernelArch", taskfile.Var{Static: m.kernelArch})
	v.Set("KernelVersion", taskfile.Var{Static: m.kernelVersion})

	t.Vars.Merge(&v)
}

func setInputVars(t *taskfile.Taskfile, inputVars []variableConfig) error {
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
