package install

import (
	"fmt"
	nrConfig "github.com/newrelic/newrelic-client-go/pkg/config"
	nrLogs "github.com/newrelic/newrelic-client-go/pkg/logs"
	"github.com/newrelic/newrelic-client-go/pkg/region"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"

	"github.com/icza/backscanner"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

var (
	assumeYes    bool
	localRecipes string
	recipeNames  []string
	recipePaths  []string
	testMode     bool
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

		logLevel := configAPI.GetLogLevel()
		config.InitFileLogger(logLevel)

		err := assertProfileIsValid()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		//TODO starting here we'd have the keys needed to create the Logs API client
		//TODO create 'temp' file - do we even need this or can we send the last x lines from the cli log?
		//tempFileName := fmt.Sprintf("%d_install_out", time.Now().UnixMilli())
		//outputFile, err := ioutil.TempFile("", tempFileName)
		//defer os.Remove(outputFile.Name())

		//TODO give RecipeInstall a new `exection.ErrorLogCollector` that is instantiated with the filename above
		i := NewRecipeInstaller(ic, client.NRClient)

		// Run the install.
		if err := i.Install(); err != nil {
			file, err := os.Open(config.GetDefaultLogFilePath())
			postMostRecentLogsToNr(100, file)

			//TODO we made it this far; this is where to prompt user, then call Logs API service (to-be-created)
			// TODO should we collect logs on interrupt?  Assuming not on update/payment required
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

func postMostRecentLogsToNr(lineCount int, logFile *os.File) {
	fileInfo, err := os.Stat(logFile.Name())

	if err != nil {
		//TODO Post this to LogsApi
		log.Debugf("Couldn't stat file: %s", logFile.Name())
		return
	}

	// building log api client
	cfg := nrConfig.New()
	cfg.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	cfg.LogLevel = "trace"
	regName, _ := region.Parse(os.Getenv("NEW_RELIC_REGION"))
	reg, _ := region.Get(regName)
	cfg.SetRegion(reg)
	cfg.Compression = nrConfig.Compression.None
	logClient := nrLogs.New(cfg)

	log.Debugf("Starting the scan")
	scanner := backscanner.New(logFile, int(fileInfo.Size()))
	currentLineCount := 0
	for {
		line, pos, err := scanner.LineBytes()
		if err != nil {
			if err == io.EOF {
				log.Debugf("Hit EOF at line position %d", pos)
			} else {
				log.Debugf("Some other error:", err)
			}
			break
		}

		if currentLineCount < lineCount {
			currentLineCount++
			logEntry := struct {
				Message string `json:"message"`
			}{
				Message: string(line),
			}

			log.Debugf("Sending log entry")
			if err := logClient.CreateLogEntry(logEntry); err != nil {
				log.Debugf("error posting Log entry: %e", err)
			} else {
				log.Infof("Just sent entry\n%v", logEntry)
			}
		}

	}
}

func assertProfileIsValid() error {
	accountID := configAPI.GetActiveProfileAccountID()
	if accountID == 0 {
		return fmt.Errorf("accountID is required")
	}

	if configAPI.GetActiveProfileString(config.APIKey) == "" {
		return fmt.Errorf("API key is required")
	}

	if configAPI.GetActiveProfileString(config.Region) == "" {
		return fmt.Errorf("region is required")
	}

	licenseKey, err := client.FetchLicenseKey(accountID, config.FlagProfileName)
	if err != nil {
		return fmt.Errorf("could not fetch license key for account %d: %s", accountID, err)
	}
	if licenseKey != configAPI.GetActiveProfileString(config.LicenseKey) {
		os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
		log.Debugf("using license key %s", utils.Obfuscate(licenseKey))
	}

	// Reinitialize client, overriding fetched values
	c, err := client.NewClient(configAPI.GetActiveProfileName())
	if err != nil {
		// An error was encountered initializing the client.  This may not be a
		// problem since many commands don't require the use of an initialized client
		log.Debugf("error initializing client: %s", err)
	}

	client.NRClient = c

	return nil
}

func init() {
	Command.Flags().StringSliceVarP(&recipePaths, "recipePath", "c", []string{}, "the path to a recipe file to install")
	Command.Flags().StringSliceVarP(&recipeNames, "recipe", "n", []string{}, "the name of a recipe to install")
	Command.Flags().BoolVarP(&testMode, "testMode", "t", false, "fakes operations for UX testing")
	Command.Flags().BoolVarP(&assumeYes, "assumeYes", "y", false, "use \"yes\" for all questions during install")
	Command.Flags().StringVarP(&localRecipes, "localRecipes", "", "", "a path to local recipes to load instead of service other fetching")
}
