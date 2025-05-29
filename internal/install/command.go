package install

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	nrRegion "github.com/newrelic/newrelic-client-go/v2/pkg/region"
)

var (
	assumeYes    bool
	localRecipes string
	recipeNames  []string
	recipePaths  []string
	tags         []string
)

// processRecipeNames validates, extracts recipe names, and sets environment variables.
func processRecipeNames(recipeNames []string) ([]string, error) {
	var extractedNames []string

	for _, recipe := range recipeNames {
		parts := strings.Split(recipe, "@")

		// Validate the recipe format
		if len(parts) < 1 || len(parts) > 2 {
			return nil, fmt.Errorf("invalid recipe format: %s", recipe)
		}

		// Extract the base recipe name
		extractedNames = append(extractedNames, parts[0])

		// If version is present, set the environment variable
		if len(parts) == 2 {
			recipeName := parts[0]
			version := parts[1]

			// Convert recipe name to uppercase, replace dashes with underscores, and append _VERSION
			envVarName := strings.ToUpper(strings.ReplaceAll(recipeName, "-", "_")) + "_VERSION"

			// Set the environment variable
			err := os.Setenv(envVarName, version)
			if err != nil {
				return nil, fmt.Errorf("error setting recipe version to environment variable %s: %v", envVarName, err)
			}
		}
	}
	return extractedNames, nil
}

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	PreRun: client.RequireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		extractedRecipeNames, err := processRecipeNames(recipeNames)
		if err != nil {
			return types.NewDetailError(types.EventTypes.OtherError, err.Error())
		}

		ic := types.InstallerContext{
			AssumeYes:    assumeYes,
			LocalRecipes: localRecipes,
			RecipeNames:  extractedRecipeNames,
			RecipePaths:  recipePaths,
		}

		ic.SetTags(tags)

		logLevel := configAPI.GetLogLevel()
		config.InitFileLogger(logLevel)

		if err := checkNetwork(); err != nil {
			return types.NewDetailError(types.EventTypes.UnableToConnect, err.Error())
		}

		detailErr := fetchLicenseKey()

		if detailErr != nil {
			log.Fatal(detailErr)
		}

		i := NewRecipeInstaller(ic)

		//// Do not install both infra and agent controls simultaneously: install only the 'agent-control' if targeted.
		//if i.IsRecipeTargeted(types.AgentControlRecipeName) && i.shouldInstallCore() {
		//	log.Debugf("'%s' is targeted, disabling infra/logs core bundle install\n", types.AgentControlRecipeName)
		//	i.shouldInstallCore = func() bool { return false }
		//}

		// Run the install.
		if err := i.Install(); err != nil {
			if err == types.ErrInterrupt {
				return nil
			}

			if _, ok := err.(*types.UpdateRequiredError); ok {
				return nil
			}

			if e, ok := err.(*nrErrors.PaymentRequiredError); ok {
				return e
			}

			if errors.Is(err, types.ErrAgentControl) {
				return err
			}

			fallbackErrorMsg := fmt.Sprintf("\nWe encountered an issue during the installation: %s.", err)
			fallbackHelpMsg := "If this problem persists, visit the documentation and support page for additional help here at https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/get-started/requirements-infrastructure-agent/"

			// In the extremely rare case we run into an uncaught error (e.g. no recipes found),
			// we need to output something to user to sinc we probably haven't displayed anything yet.
			fmt.Println(fallbackErrorMsg)
			fmt.Println(fallbackHelpMsg)
			fmt.Print("\n\n")

			log.Debug(fallbackErrorMsg)
		}

		return nil
	},
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of service other fetching")
	Command.Flags().StringSliceVarP(&tags, "tag", "", []string{}, "the tags to add during install, can be multiple. Example: --tag tag1:test,tag2:test")
}

func validateProfile() *types.DetailError {
	accountID := configAPI.GetActiveProfileAccountID()
	apiKey := configAPI.GetActiveProfileString(config.APIKey)
	region := configAPI.GetActiveProfileString(config.Region)

	if accountID == 0 {
		return types.NewDetailError(types.EventTypes.AccountIDMissing, "Account ID is required.")
	}

	if apiKey == "" {
		return types.NewDetailError(types.EventTypes.APIKeyMissing, "User API key is required.")
	}

	if !utils.IsValidUserAPIKeyFormat(apiKey) {
		return types.NewDetailError(types.EventTypes.InvalidUserAPIKeyFormat, `Invalid user API key format detected. Please provide a valid user API key. User API keys usually have a prefix of "NRAK-" or "NRAA-".`)
	}

	if region == "" {
		return types.NewDetailError(types.EventTypes.RegionMissing, "Region is required.")
	}

	if _, err := nrRegion.Parse(region); err != nil {
		return types.NewDetailError(types.EventTypes.InvalidRegion, `Invalid region provided. Valid regions are "US" or "EU".`)
	}

	return nil
}

func checkNetwork() error {
	if client.NRClient == nil {
		return nil
	}

	err := client.NRClient.TestEndpoints()

	proxyConfig := httpproxy.FromEnvironment()

	if err == nil {
		if IsProxyConfigured() {
			if strings.Contains(strings.ToLower(proxyConfig.HTTPSProxy), "http") && !strings.Contains(strings.ToLower(proxyConfig.HTTPSProxy), "https") {
				log.Warn("Please ensure the HTTPS_PROXY environment variable is set when using a proxy server.")
				log.Warn("New Relic CLI exclusively supports https proxy, not http for security reasons.")
			}
		}
	}

	if err != nil {
		if IsProxyConfigured() {
			proxyConfig := httpproxy.FromEnvironment()

			log.Warn("Proxy settings have been configured, but we are still unable to connect to the New Relic platform.")
			log.Warn("You may need to adjust your proxy environment variables or configure your proxy to allow the specified domain.")
			log.Warn("Current proxy config:")
			log.Warnf("  HTTPS_PROXY=%s", proxyConfig.HTTPSProxy)
			log.Warnf("  HTTP_PROXY=%s", proxyConfig.HTTPProxy)
			log.Warnf("  NO_PROXY=%s", proxyConfig.NoProxy)
		} else {
			log.Warn("Failed to connect to the New Relic platform.")
			log.Warn("If you need to use a proxy, consider setting the HTTPS_PROXY environment variable, then try again.")
		}

		log.Warn("More information about proxy configuration: https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy")
		log.Warn("More information about network requirements: https://docs.newrelic.com/docs/new-relic-solutions/get-started/networks/")
	}

	return err
}

// Attempt to fetch and set a license key through 3 methods:
// 1. NEW_RELIC_LICENSE_KEY environment variable,
// 2. Active profile config.LicenseKey,
// 3. API call,
// returns an error if all methods fail.
func fetchLicenseKey() *types.DetailError {
	var licenseKey string

	setLicenseKey := func(licenseKey string) {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)

		// Reinitialize client, overriding fetched values
		c, _ := client.NewClient(configAPI.GetActiveProfileName())
		client.NRClient = c

		log.Debug("using license key: ", utils.Obfuscate(licenseKey))
	}

	licenseKey = fetchLicenseKeyFromEnvironment()

	if licenseKey != "" {
		setLicenseKey(licenseKey)
		return nil
	}

	licenseKey = fetchLicenseKeyFromProfile()

	if licenseKey != "" {
		setLicenseKey(licenseKey)
		return nil
	}

	// fetch licenseKey via API
	detailErr := validateProfile()
	if detailErr != nil {
		log.Fatal(detailErr)
	}

	accountID := configAPI.GetActiveProfileAccountID()
	maxTimeoutSeconds := config.DefaultMaxTimeoutSeconds

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName, &maxTimeoutSeconds)

	if err != nil {
		err = fmt.Errorf("could not fetch license key for accountID (%d): %s",
			accountID,
			err.Error())

		return types.NewDetailError(types.EventTypes.UnableToFetchLicenseKey, err.Error())
	}

	setLicenseKey(licenseKey)
	return nil
}

func fetchLicenseKeyFromEnvironment() string {
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")

	if licenseKey == "" {
		return ""
	}

	if utils.IsValidLicenseKeyFormat(licenseKey) {
		return licenseKey
	}

	log.Debug("license key provided via NEW_RELIC_LICENSE_KEY is invalid")

	return ""
}

func fetchLicenseKeyFromProfile() string {
	// fetch licenseKey from active profile
	licenseKey := configAPI.GetActiveProfileString(config.LicenseKey)

	if licenseKey == "" {
		return ""
	}

	if utils.IsValidLicenseKeyFormat(licenseKey) {
		return licenseKey
	}

	log.Debug("license key provided by config is invalid")

	return ""
}
