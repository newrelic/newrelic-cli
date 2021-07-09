package workload

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

var (
	name                string
	entityGUIDs         []string
	entitySearchQueries []string
	scopeAccountIDs     []int
	guid                string
)

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic One workload.",
	Long: `Get a New Relic One workload

The get command retrieves a specific workload by its workload GUID.
`,
	Example: `newrelic workload get --accountId 12345678 --guid MjUyMDUyOHxOUjF8V09SS0xPQUR8MTI4Myt`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		workload, err := client.NRClient.Entities.GetEntitiesWithContext(utils.SignalCtx, []entities.EntityGUID{entities.EntityGUID(guid)})
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(workload))
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the New Relic One workloads for an account.",
	Long: `List the New Relic One workloads for an account

The list command retrieves the workloads for the given account ID.
`,
	Example: `newrelic workload list --accountId 12345678`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountIDWithFlagOverride()
		builder := entities.EntitySearchQueryBuilder{
			Type: entities.EntitySearchQueryBuilderTypeTypes.WORKLOAD,
			Tags: []entities.EntitySearchQueryBuilderTag{
				{
					Key:   "accountId",
					Value: strconv.Itoa(accountID),
				},
			},
		}
		workload, err := client.NRClient.Entities.GetEntitySearchWithContext(utils.SignalCtx, entities.EntitySearchOptions{}, "", builder, nil)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(workload))
	},
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic One workload.",
	Long: `Create a New Relic One workload

The create command accepts several different arguments for explicit and dynamic
workloads.   Multiple entity GUIDs can be provided for explicit inclusion of entities,
or multiple entity search queries can be provided for dynamic inclusion of entities.
Multiple queries will be aggregated together with an OR.  Multiple account scope
IDs can optionally be provided to include entities from different sub-accounts that
you also have access to.
`,
	Example: `newrelic workload create --name 'Example workload' --accountId 12345678 --entitySearchQuery "name like 'Example application'"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountIDWithFlagOverride()

		createInput := workloads.WorkloadCreateInput{
			Name: name,
		}

		if len(entityGUIDs) > 0 {
			for _, e := range entityGUIDs {
				createInput.EntityGUIDs = append(createInput.EntityGUIDs, entities.EntityGUID(e))
			}
		}

		if len(entitySearchQueries) > 0 {
			var queryInputs []workloads.WorkloadEntitySearchQueryInput
			for _, q := range entitySearchQueries {
				queryInputs = append(queryInputs, workloads.WorkloadEntitySearchQueryInput{Query: q})
			}
			createInput.EntitySearchQueries = queryInputs
		}

		if len(scopeAccountIDs) > 0 {
			createInput.ScopeAccounts = &workloads.WorkloadScopeAccountsInput{AccountIDs: scopeAccountIDs}
		}

		workload, err := client.NRClient.Workloads.WorkloadCreateWithContext(utils.SignalCtx, accountID, createInput)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(workload))
		log.Info("success")
	},
}

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update a New Relic One workload.",
	Long: `Update a New Relic One workload

The update command targets an existing workload by its entity GUID, and accepts
several different arguments for explicit and dynamic workloads.  Multiple entity GUIDs can
be provided for explicit inclusion of entities, or multiple entity search queries can be
provided for dynamic inclusion of entities.  Multiple queries will be aggregated
together with an OR.  Multiple account scope IDs can optionally be provided to include
entities from different sub-accounts that you also have access to.
`,
	Example: `newrelic workload update --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1' --name 'Updated workflow'`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		updateInput := workloads.WorkloadUpdateInput{
			Name: name,
		}

		if len(entityGUIDs) > 0 {
			for _, e := range entityGUIDs {
				updateInput.EntityGUIDs = append(updateInput.EntityGUIDs, entities.EntityGUID(e))
			}
		}

		if len(entitySearchQueries) > 0 {
			var queryInputs []workloads.WorkloadUpdateCollectionEntitySearchQueryInput
			for _, q := range entitySearchQueries {
				queryInputs = append(queryInputs, workloads.WorkloadUpdateCollectionEntitySearchQueryInput{Query: q})
			}
			updateInput.EntitySearchQueries = queryInputs
		}

		if len(scopeAccountIDs) > 0 {
			updateInput.ScopeAccounts = &workloads.WorkloadScopeAccountsInput{AccountIDs: scopeAccountIDs}
		}

		workload, err := client.NRClient.Workloads.WorkloadUpdateWithContext(utils.SignalCtx, entities.EntityGUID(guid), updateInput)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(workload))
		log.Info("success")
	},
}

var cmdDuplicate = &cobra.Command{
	Use:   "duplicate",
	Short: "Duplicate a New Relic One workload.",
	Long: `Duplicate a New Relic One workload

The duplicate command targets an existing workload by its entity GUID, and clones
it to the provided account ID. An optional name can be provided for the new workload.
If the name isn't specified, the name + ' copy' of the source workload is used to
compose the new name.
`,
	Example: `newrelic workload duplicate --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1' --accountId 12345678 --name 'New Workload'`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountIDWithFlagOverride()

		duplicateInput := workloads.WorkloadDuplicateInput{
			Name: name,
		}

		workload, err := client.NRClient.Workloads.WorkloadDuplicateWithContext(utils.SignalCtx, accountID, entities.EntityGUID(guid), duplicateInput)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(workload))
		log.Info("success")
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a New Relic One workload.",
	Long: `Delete a New Relic One workload

The delete command accepts a workload's entity GUID.
`,
	Example: `newrelic workload delete --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1'`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := client.NRClient.Workloads.WorkloadDeleteWithContext(utils.SignalCtx, entities.EntityGUID(guid))
		utils.LogIfFatal(err)

		log.Info("success")
	},
}

func init() {
	// Get
	Command.AddCommand(cmdGet)
	cmdGet.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload")
	utils.LogIfError(cmdGet.MarkFlagRequired("guid"))

	// List
	Command.AddCommand(cmdList)

	// Create
	Command.AddCommand(cmdCreate)
	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdCreate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "e", []string{}, "the list of entity Guids composing the workload")
	cmdCreate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdCreate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")
	utils.LogIfError(cmdCreate.MarkFlagRequired("name"))

	// Update
	Command.AddCommand(cmdUpdate)
	cmdUpdate.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload you want to update")
	cmdUpdate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdUpdate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "e", []string{}, "the list of entity Guids composing the workload")
	cmdUpdate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdUpdate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")
	utils.LogIfError(cmdUpdate.MarkFlagRequired("guid"))

	// Duplicate
	Command.AddCommand(cmdDuplicate)
	cmdDuplicate.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload you want to duplicate")
	cmdDuplicate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload to duplicate")
	utils.LogIfError(cmdDuplicate.MarkFlagRequired("guid"))

	// Delete
	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload to delete")
	utils.LogIfError(cmdDelete.MarkFlagRequired("guid"))
}
