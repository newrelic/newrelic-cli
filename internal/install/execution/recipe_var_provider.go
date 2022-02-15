package execution

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	downloadURLAccessListRegex = []string{
		`(.)?download\.newrelic\.com$`,
		`nr-downloads-ohai-(staging|testing)\.s3\.amazonaws\.com$`,
	}
)

type RecipeVarProvider struct{}

func NewRecipeVarProvider() *RecipeVarProvider {
	return &RecipeVarProvider{}
}

func (re *RecipeVarProvider) Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error) {
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

	envVarsResult := varFromEnv()

	results = append(results, systemInfoResult)
	results = append(results, profileResult)
	results = append(results, types.RecipeVariables)
	results = append(results, inputVarsResult)
	results = append(results, envVarsResult)

	for _, result := range results {
		for k, v := range result {
			vars[k] = v
		}
	}

	return vars, nil
}

func varsFromProfile(licenseKey string) (types.RecipeVars, error) {
	accountID := configAPI.GetActiveProfileString(config.AccountID)
	apiKey := configAPI.GetActiveProfileString(config.APIKey)
	region := configAPI.GetActiveProfileString(config.Region)

	if licenseKey == "" {
		return types.RecipeVars{}, errors.New("license key not found")
	}

	vars := make(types.RecipeVars)

	vars["NEW_RELIC_LICENSE_KEY"] = licenseKey
	vars["NEW_RELIC_ACCOUNT_ID"] = accountID
	vars["NEW_RELIC_API_KEY"] = apiKey
	vars["NEW_RELIC_REGION"] = region

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

func varFromEnv() types.RecipeVars {
	vars := make(types.RecipeVars)

	downloadURL := "https://download.newrelic.com/"
	envDownloadURL := os.Getenv("NEW_RELIC_DOWNLOAD_URL")
	if envDownloadURL != "" {
		URL, err := url.Parse(envDownloadURL)
		if err == nil {
			if URL.Scheme == "https" {
				for _, regexString := range downloadURLAccessListRegex {
					var regex = regexp.MustCompile(regexString)
					if regex.MatchString(URL.Host) {
						downloadURL = envDownloadURL
						break
					}
				}
			}
		} else {
			log.Warnf("Could not parse download URL: %s, detail: %s", envDownloadURL, err.Error())
		}
	}
	vars["NEW_RELIC_DOWNLOAD_URL"] = downloadURL

	vars["NEW_RELIC_CLI_LOG_FILE_PATH"] = config.GetDefaultLogFilePath()

	return vars
}
