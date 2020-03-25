package workload

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

var (
	id                  int
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

The get command retrieves a specific workload by its account ID and workload ID.
`,
	Example: `newrelic workload create --accountId 12345678 --id 1346`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			workload, err := nrClient.Workloads.GetWorkload(accountID, id)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(workload)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
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

			_, err := nrClient.Workloads.CreateWorkload(accountID, createInput)
			if err != nil {
				log.Fatal(err)
			}

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
			updateInput := workloads.UpdateInput{}

			if name != "" {
				updateInput.Name = &name
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
			if err != nil {
				log.Fatal(err)
			}

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
			var duplicateInput *workloads.DuplicateInput
			if name != "" {
				duplicateInput = &workloads.DuplicateInput{
					Name: &name,
				}
			}

			_, err := nrClient.Workloads.DuplicateWorkload(accountID, guid, duplicateInput)
			if err != nil {
				log.Fatal(err)
			}

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
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
	},
}

func init() {
	// Get
	Command.AddCommand(cmdGet)
	cmdGet.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	cmdGet.Flags().IntVarP(&id, "id", "i", 0, "the identifier of the workload")
	err := cmdGet.MarkFlagRequired("accountId")
	if err != nil {
		log.Error(err)
	}

	err = cmdGet.MarkFlagRequired("id")
	if err != nil {
		log.Error(err)
	}

	// Create
	Command.AddCommand(cmdCreate)
	cmdCreate.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdCreate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "e", []string{}, "the list of entity Guids composing the workload")
	cmdCreate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdCreate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")
	err = cmdCreate.MarkFlagRequired("accountId")
	if err != nil {
		log.Error(err)
	}

	err = cmdCreate.MarkFlagRequired("name")
	if err != nil {
		log.Error(err)
	}

	// Update
	Command.AddCommand(cmdUpdate)
	cmdUpdate.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload you want to update")
	cmdUpdate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdUpdate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "e", []string{}, "the list of entity Guids composing the workload")
	cmdUpdate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdUpdate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")

	err = cmdUpdate.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	// Duplicate
	Command.AddCommand(cmdDuplicate)
	cmdDuplicate.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload you want to duplicate")
	cmdDuplicate.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic Account ID where you want to create the new workload")
	cmdDuplicate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload to duplicate")

	err = cmdDuplicate.MarkFlagRequired("accountId")
	if err != nil {
		log.Error(err)
	}
	err = cmdDuplicate.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	// Delete
	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload to delete")
	if err != nil {
		log.Error(err)
	}
}
