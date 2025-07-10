package entities

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	eventCategory             string
	eventType                 string
	eventDescription          string
	eventEntityName           string
	eventTimestamp            int64
	eventUser                 string
	eventGroupID              string
	eventShortDescription     string
	eventCustomAttributes     string
	eventCustomAttributesFile string
	eventDataHandlingFlags    []string
)

var cmdEntityCreateEventExample = fmt.Sprintf(`newrelic entity create-event --entityName <EntityName> --category <DEPLOYMENT> --type <BASIC> --description 'desc' --timestamp %v --user 'jenkins-bot'`, time.Now().Unix())

var CmdEntityCreateEvent = &cobra.Command{
	Use:   "create-event",
	Short: "Create a New Relic change tracking event",
	Long: `Create a New Relic change tracking event.

This command allows you to create a change tracking event for a New Relic entity, including all fields supported by the API.

Custom attributes can be provided in three ways:
  1. As a JS object string via --customAttributes (e.g. '{foo: "bar", num: 2, flag: true}')
  2. As a JS object file via --customAttributesFile (e.g. --customAttributesFile ./attrs.js)
  3. From STDIN by passing --customAttributesFile - and piping the JS object (e.g. 'cat attrs.js | newrelic entity create-event ... --customAttributesFile -')

The JS object format must use unquoted keys and values of type string, boolean, or number. Example: {cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}

Validation is performed before sending to the API. Keys must be valid JS identifiers, and values must be string, boolean, or number.

For more information, see: https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/
`,
	Example: cmdEntityCreateEventExample,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		params := changetracking.ChangeTrackingCreateEventInput{}

		if eventTimestamp == 0 {
			params.Timestamp = nrtime.EpochMilliseconds(time.Now())
		} else {
			params.Timestamp = nrtime.EpochMilliseconds(time.Unix(eventTimestamp, 0))
		}

		if eventCategory == "" {
			log.Fatal("--category cannot be empty")
		}
		if eventType == "" {
			log.Fatal("--type cannot be empty")
		}
		if eventEntityName == "" {
			log.Fatal("--entityName cannot be empty")
		}

		params.Description = eventDescription
		params.User = eventUser
		params.GroupId = eventGroupID
		params.ShortDescription = eventShortDescription
		params.EntitySearch = changetracking.ChangeTrackingEntitySearchInput{
			Query: fmt.Sprintf("name = '%s'", eventEntityName),
		}
		params.CategoryAndTypeData = &changetracking.ChangeTrackingCategoryRelatedInput{
			Kind: &changetracking.ChangeTrackingCategoryAndTypeInput{
				Category: eventCategory,
				Type:     eventType,
			},
		}

		// Custom Attributes: support --customAttributes, --customAttributesFile, and STDIN
		var customAttrRaw string
		if eventCustomAttributesFile != "" {
			if eventCustomAttributesFile == "-" {
				// Read from STDIN
				stdinBytes, err := os.ReadFile("/dev/stdin")
				if err != nil {
					log.Fatalf("Failed to read custom attributes from STDIN: %v", err)
				}
				customAttrRaw = string(stdinBytes)
			} else {
				fileBytes, err := os.ReadFile(eventCustomAttributesFile)
				if err != nil {
					log.Fatalf("Failed to read custom attributes file: %v", err)
				}
				customAttrRaw = string(fileBytes)
			}
		} else if eventCustomAttributes != "" {
			customAttrRaw = eventCustomAttributes
		}

		if customAttrRaw != "" {
			// Validate JS object format: keys must be valid JS identifiers, values must be string, bool, or number
			// Accepts: {foo: "bar", num: 2, flag: true}
			// This is a basic validation, not a full JS parser
			jsObj := strings.TrimSpace(customAttrRaw)
			if !strings.HasPrefix(jsObj, "{") || !strings.HasSuffix(jsObj, "}") {
				log.Fatal("customAttributes must be a JS object, e.g. {foo: \"bar\", num: 2, flag: true}")
			}
			// Validate keys and values (simple regex)
			kvRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*([^"]+|"[^"]*"|true|false|[0-9.]+)`) // key: value
			matches := kvRe.FindAllStringSubmatch(jsObj, -1)
			if len(matches) == 0 {
				log.Fatal("customAttributes must contain at least one valid key: value pair")
			}
			// Optionally, further validation can be added here
			// Convert JS object string to map[string]interface{} for API
			attrs, err := changetracking.ReadCustomAttributesJS(customAttrRaw, false)
			if err != nil {
				log.Fatalf("Failed to parse customAttributes as JS object: %v", err)
			}
			params.CustomAttributes = attrs
		}

		// Parse data handling flags
		var flags []changetracking.ChangeTrackingValidationFlag
		for _, flag := range eventDataHandlingFlags {
			switch flag {
			case "ALLOW_CUSTOM_CATEGORY_OR_TYPE":
				flags = append(flags, changetracking.ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE)
			case "FAIL_ON_FIELD_LENGTH":
				flags = append(flags, changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH)
			case "FAIL_ON_REST_API_FAILURES":
				flags = append(flags, changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_REST_API_FAILURES)
			}
		}
		dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{ValidationFlags: flags}

		result, err := client.NRClient.ChangeTracking.ChangeTrackingCreateEventWithContext(
			utils.SignalCtx,
			params,
			dataHandlingRules,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdEntityCreateEvent)
	CmdEntityCreateEvent.Flags().StringVar(&eventEntityName, "entityName", "", "the name of the entity associated with this event. entityName is required.")
	utils.LogIfError(CmdEntityCreateEvent.MarkFlagRequired("entityName"))

	CmdEntityCreateEvent.Flags().StringVar(&eventCategory, "category", "", "category of event, e.g., DEPLOYMENT, CONFIG_CHANGE, etc. category is required.")
	utils.LogIfError(CmdEntityCreateEvent.MarkFlagRequired("category"))

	CmdEntityCreateEvent.Flags().StringVar(&eventType, "type", "", "type of event, e.g., BASIC, ROLLBACK, etc. type is required.")
	utils.LogIfError(CmdEntityCreateEvent.MarkFlagRequired("type"))

	CmdEntityCreateEvent.Flags().StringVar(&eventDescription, "description", "", "a description of the event")
	CmdEntityCreateEvent.Flags().StringVar(&eventUser, "user", "", "username of the actor or bot")
	CmdEntityCreateEvent.Flags().StringVar(&eventGroupID, "groupId", "", "string that can be used to correlate two or more events")
	CmdEntityCreateEvent.Flags().StringVar(&eventShortDescription, "shortDescription", "", "short description for the event")
	CmdEntityCreateEvent.Flags().StringVar(&eventCustomAttributes, "customAttributes", "", "custom attributes for the event in JS object format, e.g. {key1: 'value1', key2: 2}")
	CmdEntityCreateEvent.Flags().StringVar(&eventCustomAttributesFile, "customAttributesFile", "", "path to a file containing custom attributes in JS object format, or '-' to read from STDIN")
	CmdEntityCreateEvent.Flags().StringSliceVar(&eventDataHandlingFlags, "dataHandlingFlags", []string{"FAIL_ON_FIELD_LENGTH"}, "array of data handling flags, e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE,FAIL_ON_FIELD_LENGTH,FAIL_ON_REST_API_FAILURES")
	CmdEntityCreateEvent.Flags().Int64VarP(&eventTimestamp, "timestamp", "t", 0, "the time of the event, the number of milliseconds since the Unix epoch, defaults to now")
}
