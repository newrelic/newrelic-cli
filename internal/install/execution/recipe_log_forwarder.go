package execution

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	nrConfig "github.com/newrelic/newrelic-client-go/v2/pkg/config"
	nrLogs "github.com/newrelic/newrelic-client-go/v2/pkg/logs"
	"github.com/newrelic/newrelic-client-go/v2/pkg/region"
)

type LogForwarder interface {
	PromptUserToSendLogs(reader io.Reader) bool
	SendLogsToNewRelic(recipeName string, data []string)
	HasUserOptedIn() bool
	SetUserOptedIn(val bool)
}

type LogEntry struct {
	Attributes map[string]interface{} `json:"attributes"`
	LogType    string                 `json:"logType"`
	Message    string                 `json:"message"`
}

type RecipeLogForwarder struct {
	LogEntries []LogEntry
	optIn      bool
}

func NewRecipeLogForwarder() *RecipeLogForwarder {
	return &RecipeLogForwarder{
		LogEntries: []LogEntry{},
		optIn:      false,
	}
}

func (lf *RecipeLogForwarder) PromptUserToSendLogs(reader io.Reader) bool {
	fmt.Printf("\n%s Something went wrong during the installation — let’s look at the docs and try again. Would you like to send us the logs to help us understand what happened? [Y/n] ", color.YellowString("\u0021"))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		input := strings.TrimSuffix(scanner.Text(), "\n")
		if strings.ToUpper(input) == "Y" || len(strings.TrimSpace(input)) == 0 {
			return true
		}
		return false
	}
	if scanner.Err() != nil {
		log.Debugf("An error occurred while reading input: %e", scanner.Err())
	}

	return false
}

func (lf *RecipeLogForwarder) HasUserOptedIn() bool {
	return lf.optIn
}

func (lf *RecipeLogForwarder) SetUserOptedIn(val bool) {
	lf.optIn = val
}

func (lf *RecipeLogForwarder) SendLogsToNewRelic(recipeName string, data []string) {
	lf.buildLogEntryBatch(recipeName, data)

	// building log api client
	config, err := createLogClientConfig()
	if nil != err {
		log.Debugf("Could not configure New Relic LogsApi client: %e", err)
		return
	}
	logClient := nrLogs.New(config)

	// fetch accountId & configure client for batch mode
	accountID, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
	if nil != err {
		log.Debugf("Could not determine account id for log destintation: %e", err)
		return
	}
	if err := logClient.BatchMode(context.Background(), accountID); err != nil {
		log.Fatal("error starting batch mode: ", err)
	}

	// enqueue log entries.
	for _, logEntry := range lf.LogEntries {
		if err := logClient.EnqueueLogEntry(context.Background(), logEntry); err != nil {
			log.Fatal("error queuing log entry: ", err)
		}
	}

	// Force flush/send; sleep seems necessary, otherwise logs don't appear to land in NR
	time.Sleep(5 * time.Second)
	if err := logClient.Flush(); err != nil {
		log.Fatal("error flushing log queue: ", err)
	}
}

func (lf *RecipeLogForwarder) buildLogEntryBatch(recipeName string, data []string) {
	now := time.Now().UnixMilli()
	for _, line := range data {
		now++ //using timestamp to retain log sequence
		lf.LogEntries = append(lf.LogEntries, LogEntry{map[string]interface{}{"nr-install-recipe": recipeName, "timestamp": now}, "cli-output", line})
	}
}

func createLogClientConfig() (nrConfig.Config, error) {
	cfg := nrConfig.Config{
		LicenseKey:  os.Getenv("NEW_RELIC_LICENSE_KEY"),
		LogLevel:    "info",
		Compression: nrConfig.Compression.None,
	}

	regName, _ := region.Parse(os.Getenv("NEW_RELIC_REGION"))
	reg, _ := region.Get(regName)
	err := cfg.SetRegion(reg)
	if nil != err {
		log.Debugf("Could not set region on LogsApi client: %e", err)
		return cfg, err
	}

	return cfg, nil
}
