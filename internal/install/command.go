package install

import (
	"errors"
	"fmt"
	"os"

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
	testMode     bool
	tags         []string
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	PreRun: client.RequireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		ic := types.InstallerContext{
			AssumeYes:    assumeYes,
			LocalRecipes: localRecipes,
			RecipeNames:  recipeNames,
			RecipePaths:  recipePaths,
		}
		ic.SetTags(tags)

		logLevel := configAPI.GetLogLevel()
		config.InitFileLogger(logLevel)

		detailErr := validateProfile()
		if detailErr != nil {
			log.Fatal(detailErr)
		}

		detailErr = fetchLicenseKey()
		if detailErr != nil {
			log.Fatal(detailErr)
		}

		// Reinitialize client, overriding fetched values
		c, _ := client.NewClient(configAPI.GetActiveProfileName())
		client.NRClient = c

		i := NewRecipeInstaller(ic, c)

		//// Do not install both infra and super agents simultaneously: install only the 'super-agent' if targeted.
		//if i.IsRecipeTargeted(types.SuperAgentRecipeName) && i.shouldInstallCore() {
		//	log.Debugf("'%s' is targeted, disabling infra/logs core bundle install\n", types.SuperAgentRecipeName)
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

			if errors.Is(err, types.ErrSuperAgent) {
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
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
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

	if err := checkNetwork(); err != nil {
		return types.NewDetailError(types.EventTypes.UnableToConnect, err.Error())
	}

	return nil
}

func checkNetwork() error {
	if client.NRClient == nil {
		return nil
	}

	err := client.NRClient.TestEndpoints()
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
	accountID := configAPI.GetActiveProfileAccountID()

	var licenseKey string

	defer func() {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
		log.Debug("using license key: ", utils.Obfuscate(licenseKey))
	}()

	// fetch licenseKey from environment
	licenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")

	if utils.IsValidLicenseKeyFormat(licenseKey) {
		return nil
	} else {
		log.Debug("license key provided via NEW_RELIC_LICENSE_KEY is invalid")
	}

	// fetch licenseKey from active profile
	licenseKey = configAPI.GetActiveProfileString(config.LicenseKey)

	if utils.IsValidLicenseKeyFormat(licenseKey) {
		return nil
	} else {
		log.Debug("license key provided by config is invalid")
	}

	// fetch licenseKey via API
	maxTimeoutSeconds := config.DefaultMaxTimeoutSeconds

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName, &maxTimeoutSeconds)

	if err != nil {
		err = fmt.Errorf("could not fetch license key for accountID (%d): %s",
			accountID,
			err.Error())

		return types.NewDetailError(types.EventTypes.UnableToFetchLicenseKey, err.Error())
	}

	return nil
}
