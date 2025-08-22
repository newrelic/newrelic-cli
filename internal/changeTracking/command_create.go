package changeTracking

import (
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
	eventCategory         string
	eventType             string
	eventDescription      string
	eventSearchQuery      string
	eventTimestamp        int64
	eventUser             string
	eventGroupID          string
	eventShortDescription string
	eventCustomAttributes string
	eventValidationFlags  []string

	// Deployment fields
	eventVersion   string
	eventChangelog string
	eventCommit    string
	eventDeepLink  string

	// Feature flag fields
	eventFeatureFlagId string
)

var cmdChangeTrackingCreateExample = `# Deployment event example with complex entity search (shows nested quotes pattern)
newrelic changeTracking create \
  --entitySearch "name = 'MyService' AND type = 'SERVICE'" \
  --category Deployment \
  --type Basic \
  --description "Deployed version 1.2.3 to production" \
  --version "1.2.3" \
  --changelog "https://github.com/myorg/myservice/releases/tag/v1.2.3" \
  --commit "abc123def456" \
  --user "ci-cd-bot"

# Feature Flag event example (shows both nested quote patterns)
newrelic changeTracking create \
  --entitySearch "name = 'MyApp'" \
  --category "Feature Flag" \
  --type Basic \
  --description "Enabled new checkout flow feature flag" \
  --featureFlagId "checkout-flow-v2" \
  --validationFlags FAIL_ON_FIELD_LENGTH \
  --user "product-team"

# Custom category/type example with custom attributes and multiple validation flags (comma-separated style)
newrelic changeTracking create \
  --entitySearch "id = '<Entity GUID>'" \
  --category Operational \
  --type "Maintenance Window" \
  --description "Database maintenance completed" \
  --customAttributes '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}' \
  --validationFlags ALLOW_CUSTOM_CATEGORY_OR_TYPE,FAIL_ON_FIELD_LENGTH,FAIL_ON_REST_API_FAILURES \
  --user "ops-team"

# Custom category/type example using repeated --validationFlags style (each flag separately)
newrelic changeTracking create \
  --entitySearch "id = '<Entity GUID>'" \
  --category Operational \
  --type "Maintenance Window" \
  --description "Database maintenance completed" \
  --customAttributes '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}' \
  --validationFlags ALLOW_CUSTOM_CATEGORY_OR_TYPE \
  --validationFlags FAIL_ON_FIELD_LENGTH \
  --validationFlags FAIL_ON_REST_API_FAILURES \
  --user "ops-team"`

var CmdChangeTrackingCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic change tracking event",
	Long: `Create a New Relic change tracking event.

This command allows you to create a change tracking event for a New Relic entity, supporting all fields in the Change Tracking GraphQL API schema for the changeTrackingCreateEvent mutation. For more information on each field, visit: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#change-tracking-event-mutation

Required fields:
  --entitySearch        Entity search query (e.g., name = 'MyService' AND type = 'SERVICE'). See our docs on 'entitySearch.query' under 'Required attributes' for more detailed examples: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#required-fields
  --category            Category of event (e.g. Deployment, Feature Flag, Operational, etc.)
  --type                Type of event (e.g. Basic, Rollback, Server Reboot, etc.)

For Deployment events, the following are required/supported:
  --version             Version of the deployment (required)
  --changelog           Changelog for the deployment (URL or text)
  --commit              Commit hash for the deployment
  --deepLink            Deep link URL for the deployment

For Feature Flag events, the following are required/supported:
  --featureFlagId       ID of the feature flag (required)

Other supported fields:
  --description         Description of the event
  --user                Username of the actor or bot
  --groupId             String to correlate two or more events
  --shortDescription    Short description for the event
  --customAttributes    Custom attributes: use '-' for STDIN, '{...}' for inline JS object, or provide a file path
  --validationFlags     Comma-separated list of validation flags (e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE, FAIL_ON_FIELD_LENGTH, FAIL_ON_REST_API_FAILURES)
  --timestamp           Time of the event (milliseconds since Unix epoch, defaults to now). Can not be more than 24 hours in the past or future

Custom attributes can be provided in three ways:
  1. From STDIN by passing '-' (e.g. 'echo  '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}' | newrelic changeTracking create ... --customAttributes -')
  2. As an inline JS object starting with '{' (e.g. --customAttributes '{cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}')
  3. As a file path (e.g. --customAttributes ./attrs.js)

The JS object format must use unquoted keys and values of type string, boolean, or number. Example: {cloud_vendor: "vendor_name", region: "us-east-1", isProd: true, instances: 2}

Validation is performed before sending to the API. Keys must be valid JS identifiers, and values must be string, boolean, or number.

For more information, see: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/#change-tracking-event-mutation
`,
	Example: cmdChangeTrackingCreateExample,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		params := changetracking.ChangeTrackingCreateEventInput{}

		if eventTimestamp == 0 {
			params.Timestamp = nrtime.EpochMilliseconds(time.Now())
		} else {
			params.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(eventTimestamp))
		}

		if eventCategory == "" {
			log.Fatal("--category cannot be empty")
		}
		if eventType == "" {
			log.Fatal("--type cannot be empty")
		}
		if eventSearchQuery == "" {
			log.Fatal("--entitySearch cannot be empty")
		}

		params.Description = eventDescription
		params.User = eventUser
		params.GroupId = eventGroupID
		params.ShortDescription = eventShortDescription
		params.EntitySearch = changetracking.ChangeTrackingEntitySearchInput{
			Query: eventSearchQuery,
		}
		params.CategoryAndTypeData = &changetracking.ChangeTrackingCategoryRelatedInput{
			Kind: &changetracking.ChangeTrackingCategoryAndTypeInput{
				Category: eventCategory,
				Type:     eventType,
			},
			CategoryFields: &changetracking.ChangeTrackingCategoryFieldsInput{},
		}
		// Set deployment fields if category is DEPLOYMENT
		if strings.ToUpper(eventCategory) == "DEPLOYMENT" {
			if eventVersion == "" {
				log.Fatal("--version is required for DEPLOYMENT events")
			}
			params.CategoryAndTypeData.CategoryFields.Deployment = &changetracking.ChangeTrackingDeploymentFieldsInput{
				Version:   eventVersion,
				Changelog: eventChangelog,
				Commit:    eventCommit,
				DeepLink:  eventDeepLink,
			}
		}
		// Set feature flag fields if category is FEATURE FLAG
		if strings.ToUpper(eventCategory) == "FEATURE FLAG" {
			if eventFeatureFlagId == "" {
				log.Fatal("--featureFlagId is required for FEATURE FLAG events")
			}
			params.CategoryAndTypeData.CategoryFields.FeatureFlag = &changetracking.ChangeTrackingFeatureFlagFieldsInput{
				FeatureFlagId: eventFeatureFlagId,
			}
		}

		// Custom Attributes: support --customAttributes with three parsing modes:
		// 1. If equals "-", read from STDIN
		// 2. If starts with "{", parse as JS object
		// 3. Otherwise, treat as file path
		var customAttrRaw string
		if eventCustomAttributes != "" {
			if eventCustomAttributes == "-" {
				// Read from STDIN
				stdinBytes, err := os.ReadFile("/dev/stdin")
				if err != nil {
					log.Fatalf("Failed to read custom attributes from STDIN: %v", err)
				}
				customAttrRaw = string(stdinBytes)
			} else if strings.HasPrefix(strings.TrimSpace(eventCustomAttributes), "{") {
				// Parse as JS object directly
				customAttrRaw = eventCustomAttributes
			} else {
				// Treat as file path
				fileBytes, err := os.ReadFile(eventCustomAttributes)
				if err != nil {
					log.Fatalf("Failed to read custom attributes file: %v", err)
				}
				customAttrRaw = string(fileBytes)
			}
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
			kvRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*([^\"]+|\"[^\"]*\"|true|false|[0-9.]+)`) // key: value
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

		// Parse validation flags
		var flags []changetracking.ChangeTrackingValidationFlag
		for _, flag := range eventValidationFlags {
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
	Command.AddCommand(CmdChangeTrackingCreate)
	CmdChangeTrackingCreate.Flags().StringVar(&eventSearchQuery, "entitySearch", "", "the NRQL entity search query for this event. Example: name = 'MyService' AND type = 'SERVICE' (required)")
	utils.LogIfError(CmdChangeTrackingCreate.MarkFlagRequired("entitySearch"))

	CmdChangeTrackingCreate.Flags().StringVar(&eventCategory, "category", "", "category of event, e.g., DEPLOYMENT, CONFIG_CHANGE, etc. category is required.")
	utils.LogIfError(CmdChangeTrackingCreate.MarkFlagRequired("category"))

	CmdChangeTrackingCreate.Flags().StringVar(&eventType, "type", "", "type of event, e.g., BASIC, ROLLBACK, etc. type is required.")
	utils.LogIfError(CmdChangeTrackingCreate.MarkFlagRequired("type"))

	CmdChangeTrackingCreate.Flags().StringVar(&eventDescription, "description", "", "a description of the event")
	CmdChangeTrackingCreate.Flags().StringVar(&eventUser, "user", "", "username of the actor or bot")
	CmdChangeTrackingCreate.Flags().StringVar(&eventGroupID, "groupId", "", "string that can be used to correlate two or more events")
	CmdChangeTrackingCreate.Flags().StringVar(&eventShortDescription, "shortDescription", "", "short description for the event")
	CmdChangeTrackingCreate.Flags().StringVar(&eventCustomAttributes, "customAttributes", "", "custom attributes: use '-' for STDIN, '{...}' for inline JS object, or provide a file path")
	CmdChangeTrackingCreate.Flags().StringSliceVar(&eventValidationFlags, "validationFlags", nil, "comma-separated list of validation flags, e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE,FAIL_ON_FIELD_LENGTH,FAIL_ON_REST_API_FAILURES")
	CmdChangeTrackingCreate.Flags().Int64VarP(&eventTimestamp, "timestamp", "t", 0, "the time of the event, the number of milliseconds since the Unix epoch, defaults to now")

	// Deployment fields
	CmdChangeTrackingCreate.Flags().StringVar(&eventVersion, "version", "", "version of the deployment")
	CmdChangeTrackingCreate.Flags().StringVar(&eventChangelog, "changelog", "", "changelog for the deployment")
	CmdChangeTrackingCreate.Flags().StringVar(&eventCommit, "commit", "", "commit hash for the deployment")
	CmdChangeTrackingCreate.Flags().StringVar(&eventDeepLink, "deepLink", "", "deep link URL for the deployment")

	// Feature flag fields
	CmdChangeTrackingCreate.Flags().StringVar(&eventFeatureFlagId, "featureFlagId", "", "ID of the feature flag")
}
