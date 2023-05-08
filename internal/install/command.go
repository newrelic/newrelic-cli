package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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

		err := assertProfileIsValid(config.DefaultMaxTimeoutSeconds, sg)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		// Reinitialize client, overriding fetched values
		c, err := client.NewClient(configAPI.GetActiveProfileName())
		if err != nil {
			// An error was encountered initializing the client.  This may not be a
			// problem since many commands don't require the use of an initialized client
			log.Debugf("error initializing client: %s", err)
			sg.TrackInfo(segment.EventTypes.UnableToOverrideClient, segment.NewEventInfo(err.Error()))
		}

		client.NRClient = c

		i := NewRecipeInstaller(ic, client.NRClient)

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

func assertProfileIsValid(maxTimeoutSeconds int, sg *segment.Segment) error {

	accountID := configAPI.GetActiveProfileAccountID()
	sg.Track(segment.EventTypes.InstallStarted)

	if accountID == 0 {
		sg.Track(segment.EventTypes.AccountIDMissing)
		return fmt.Errorf("accountID is required")
	}

	if configAPI.GetActiveProfileString(config.APIKey) == "" {
		sg.Track(segment.EventTypes.APIKeyMissing)
		return fmt.Errorf("API key is required")
	}

	if configAPI.GetActiveProfileString(config.Region) == "" {
		sg.Track(segment.EventTypes.RegionMissing)
		return fmt.Errorf("region is required")
	}

	if err := checkNetwork(client.NRClient); err != nil {
		sg.Track(segment.EventTypes.UnableToConnect)
		return err
	}

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName, &maxTimeoutSeconds)
	if err != nil {
		sg.TrackInfo(segment.EventTypes.UnableToFetchLicenseKey, segment.NewEventInfo(err.Error()))
		return fmt.Errorf("could not fetch license key for account %d:, license key: %v %s", accountID, utils.Obfuscate(licenseKey), err)
	}
	sg.Track(segment.EventTypes.LicenseKeyFetchedOk)

	if licenseKey != configAPI.GetActiveProfileString(config.LicenseKey) {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
		log.Debugf("using license key %s", utils.Obfuscate(licenseKey))
	}

	return nil
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of service other fetching")
	Command.Flags().StringSliceVarP(&tags, "tag", "", []string{}, "the tags to add during install, can be multiple. Example: --tag tag1:test,tag2:test")
}
