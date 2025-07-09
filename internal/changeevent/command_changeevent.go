package changeevent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	eventCategory          string
	eventType              string
	eventDescription       string
	eventEntityName        string
	eventTimestamp         int64
	eventUser              string
	eventGroupID           string
	eventShortDescription  string
	eventCustomAttributes  string
	eventDataHandlingFlags []string
)

var CmdChangeTracking = &cobra.Command{
	Use:   "changetracking",
	Short: "Manage change tracking events for New Relic",
}

var cmdChangeTrackingCreateEventExample = fmt.Sprintf(`newrelic changetracking create-event --entityName <EntityName> --category <DEPLOYMENT> --type <BASIC> --description 'desc' --timestamp %v --user 'jenkins-bot'`, time.Now().Unix())

var CmdChangeTrackingCreateEvent = &cobra.Command{
	Use:   "create-event",
	Short: "Create a New Relic change tracking event",
	Long: `Create a New Relic change tracking event

The create-event command marks a change for a New Relic entity.
	`,
	Example: cmdChangeTrackingCreateEventExample,
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
		if eventCustomAttributes != "" {
			var customAttrs changetracking.ChangeTrackingRawCustomAttributesMap
			if err := json.Unmarshal([]byte(eventCustomAttributes), &customAttrs); err != nil {
				log.Fatalf("Invalid customAttributes JSON: %v", err)
			}
			params.CustomAttributes = customAttrs
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
	CmdChangeTracking.AddCommand(CmdChangeTrackingCreateEvent)

	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventEntityName, "entityName", "", "the name of the entity associated with this event. entityName is required.")
	utils.LogIfError(CmdChangeTrackingCreateEvent.MarkFlagRequired("entityName"))

	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventCategory, "category", "", "category of event, e.g., DEPLOYMENT, CONFIG_CHANGE, etc. category is required.")
	utils.LogIfError(CmdChangeTrackingCreateEvent.MarkFlagRequired("category"))

	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventType, "type", "", "type of event, e.g., BASIC, ROLLBACK, etc. type is required.")
	utils.LogIfError(CmdChangeTrackingCreateEvent.MarkFlagRequired("type"))

	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventDescription, "description", "", "a description of the event")
	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventUser, "user", "", "username of the actor or bot")
	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventGroupID, "groupId", "", "string that can be used to correlate two or more events")
	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventShortDescription, "shortDescription", "", "short description for the event")
	CmdChangeTrackingCreateEvent.Flags().StringVar(&eventCustomAttributes, "customAttributes", "", "custom attributes for the event in JSON object format, e.g. {key1: 'value1', key2: 2}")
	CmdChangeTrackingCreateEvent.Flags().StringSliceVar(&eventDataHandlingFlags, "dataHandlingFlags", []string{"FAIL_ON_FIELD_LENGTH"}, "array of data handling flags, e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE,FAIL_ON_FIELD_LENGTH,FAIL_ON_REST_API_FAILURES")
	CmdChangeTrackingCreateEvent.Flags().Int64VarP(&eventTimestamp, "timestamp", "t", 0, "the time of the event, the number of milliseconds since the Unix epoch, defaults to now")
}
