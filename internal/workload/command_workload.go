package workload

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

var (
	accountID           int
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

The get command retrieves a specific workload by its account ID and workload GUID.
`,
	Example: `newrelic workload create --accountId 12345678 --guid MjUyMDUyOHxOUjF8V09SS0xPQUR8MTI4Myt`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			workload, err := nrClient.Workloads.GetWorkload(accountID, guid)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(workload))
		})
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the New Relic One workloads for an account.",
	Long: `List the New Relic One workloads for an account

The list command retrieves the workloads for the given account ID.
`,
	Example: `newrelic workload list --accountId 12345678`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			workload, err := nrClient.Workloads.ListWorkloads(accountID)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(workload))
		})
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
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			createInput := workloads.CreateInput{
				Name: name,
			}

			if len(entityGUIDs) > 0 {
				createInput.EntityGUIDs = entityGUIDs
			}

			if len(entitySearchQueries) > 0 {
				var queryInputs []workloads.EntitySearchQueryInput
				for _, q := range entitySearchQueries {
					queryInputs = append(queryInputs, workloads.EntitySearchQueryInput{Query: q})
				}
				createInput.EntitySearchQueries = queryInputs
			}

			if len(scopeAccountIDs) > 0 {
				createInput.ScopeAccountsInput = &workloads.ScopeAccountsInput{AccountIDs: scopeAccountIDs}
			}

			workload, err := nrClient.Workloads.CreateWorkload(accountID, createInput)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(workload))
			log.Info("success")
		})
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
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			updateInput := workloads.UpdateInput{
				Name: name,
			}

			if len(entityGUIDs) > 0 {
				updateInput.EntityGUIDs = entityGUIDs
			}

			if len(entitySearchQueries) > 0 {
				var queryInputs []workloads.EntitySearchQueryInput
				for _, q := range entitySearchQueries {
					queryInputs = append(queryInputs, workloads.EntitySearchQueryInput{Query: q})
				}
				updateInput.EntitySearchQueries = queryInputs
			}

			if len(scopeAccountIDs) > 0 {
				updateInput.ScopeAccountsInput = &workloads.ScopeAccountsInput{AccountIDs: scopeAccountIDs}
			}

			_, err := nrClient.Workloads.UpdateWorkload(guid, updateInput)
			utils.LogIfFatal(err)

			log.Info("success")
		})
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
	Example: `newrelic workload duplicate --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1' --accountID 12345678 --name 'New Workload'`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			duplicateInput := &workloads.DuplicateInput{
				Name: name,
			}

			workload, err := nrClient.Workloads.DuplicateWorkload(accountID, guid, duplicateInput)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(workload))
			log.Info("success")
		})
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a New Relic One workload.",
	Long: `Delete a New Relic One workload

The delete command accepts a workload's entity GUID.
`,
	Example: `newrelic workload delete --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1'`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			_, err := nrClient.Workloads.DeleteWorkload(guid)
			utils.LogIfFatal(err)

			log.Info("success")
		})
	},
}

func init() {
	// Get
	Command.AddCommand(cmdGet)
	cmdGet.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	cmdGet.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload")
	utils.LogIfError(cmdGet.MarkFlagRequired("accountId"))
	utils.LogIfError(cmdGet.MarkFlagRequired("id"))

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	utils.LogIfError(cmdList.MarkFlagRequired("accountId"))

	// Create
	Command.AddCommand(cmdCreate)
	cmdCreate.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdCreate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "e", []string{}, "the list of entity Guids composing the workload")
	cmdCreate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdCreate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")
	utils.LogIfError(cmdCreate.MarkFlagRequired("accountId"))
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
	cmdDuplicate.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic Account ID where you want to create the new workload")
	cmdDuplicate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload to duplicate")
	utils.LogIfError(cmdDuplicate.MarkFlagRequired("accountId"))
	utils.LogIfError(cmdDuplicate.MarkFlagRequired("guid"))

	// Delete
	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload to delete")
	utils.LogIfError(cmdDelete.MarkFlagRequired("guid"))
}
