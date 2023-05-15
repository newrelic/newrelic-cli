package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/segment"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
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

		sg := initSegment()
		sg.Track(types.EventTypes.InstallStarted)
		detailErr := validateProfile(config.DefaultMaxTimeoutSeconds)

		if detailErr != nil {
			ei := segment.NewEventInfo(detailErr.EventName, detailErr.Details)
			sg.TrackInfo(ei)
			sg.Close()
			log.Fatal(detailErr)
		}

		// Reinitialize client, overriding fetched values
		c, _ := client.NewClient(configAPI.GetActiveProfileName())
		client.NRClient = c

		i := NewRecipeInstaller(ic, c, sg)

		// Run the install.
		if err := i.Install(); err != nil {
			defer sg.Close()
			if err == types.ErrInterrupt {
				return nil
			}

			if _, ok := err.(*types.UpdateRequiredError); ok {
				return nil
			}

			if e, ok := err.(*nrErrors.PaymentRequiredError); ok {
				return e
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

func initSegment() *segment.Segment {
	accountID := configAPI.GetActiveProfileAccountID()
	region := configAPI.GetActiveProfileString(config.Region)
	isProxyConfigured := IsProxyConfigured()
	writeKey, err := recipes.NewEmbeddedRecipeFetcher().GetSegmentWriteKey()
	if err != nil {
		log.Debug("segment: error reading write key, cannot write to segment", err)
		return nil
	}

	return segment.New(writeKey, accountID, region, isProxyConfigured)
}

func validateProfile(maxTimeoutSeconds int) *types.DetailError {
	accountID := configAPI.GetActiveProfileAccountID()
	APIKey := configAPI.GetActiveProfileString(config.APIKey)
	region := configAPI.GetActiveProfileString(config.Region)

	if accountID == 0 {
		return types.NewDetailError(types.EventTypes.AccountIDMissing, "account ID is required")
	}

	if APIKey == "" {
		return types.NewDetailError(types.EventTypes.APIKeyMissing, "API key is required")
	}

	if region == "" {
		return types.NewDetailError(types.EventTypes.RegionMissing, "region is required")
	}

	if err := checkNetwork(); err != nil {
		return types.NewDetailError(types.EventTypes.UnableToConnect, err.Error())
	}

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName, &maxTimeoutSeconds)
	if err != nil {
		details := fmt.Sprintf("could not fetch license key for account %d:, license key: %v %s", accountID, utils.Obfuscate(licenseKey), err)
		return types.NewDetailError(types.EventTypes.UnableToFetchLicenseKey, details)
	}

	if licenseKey != configAPI.GetActiveProfileString(config.LicenseKey) {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
		log.Debugf("using license key %s", utils.Obfuscate(licenseKey))
	}

	return nil
}

func checkNetwork() error {
	err := client.NRClient.TestEndpoints()
	if err != nil {
		if IsProxyConfigured() {
			log.Warn("Proxy settings have been configured, but we are still unable to connect to the New Relic platform.")
			log.Warn("You may need to adjust your proxy environment variables or configure your proxy to allow the specified domain.")
			log.Warn("Current proxy config:")
			proxyConfig := httpproxy.FromEnvironment()
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
