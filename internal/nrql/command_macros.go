package nrql

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

var (
	macroID         string
	macroName       string
	macroExpression string
	macroOrgID      string
	macroNumParams  int
)

var cmdMacros = &cobra.Command{
	Use:   "macros",
	Short: "Commands for managing NRQL Macros",
	Long: `Commands for managing NRQL Macros

NRQL Macros are reusable NRQL expression fragments that can be referenced
as $macroName in queries. They are scoped to an organization.
`,
	Example: "newrelic nrql macros --help",
}

var cmdMacrosCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a NRQL Macro",
	Long: `Create a NRQL Macro

Creates a new NRQL Macro scoped to your organization. The macro name is
immutable after creation. The expression is a NRQL fragment that will be
substituted wherever $macroName appears in a query.
`,
	Example: "newrelic nrql macros create --name filterProd --expression \"environment = 'production'\"",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		orgID := macroOrgID
		if orgID == "" {
			org, err := client.NRClient.Organization.GetOrganization()
			if err != nil {
				log.Fatal("could not determine organization ID: ", err)
			}
			orgID = org.ID
		}

		result, err := client.NRClient.Nrqlmacros.EntityManagementCreateNRQLMacro(
			nrqlmacros.EntityManagementNRQLMacroEntityCreateInput{
				Name:       macroName,
				Expression: macroExpression,
				NumParameters: macroNumParams,
				Scope: nrqlmacros.EntityManagementScopedReferenceInput{
					Type: nrqlmacros.EntityManagementEntityScopeTypes.ORGANIZATION,
					ID:   orgID,
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Print(result)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosGet = &cobra.Command{
	Use:   "get",
	Short: "Get a NRQL Macro by ID",
	Long: `Get a NRQL Macro by ID

Retrieves a single NRQL Macro by its entity ID.
`,
	Example: "newrelic nrql macros get --id <macroID>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		entityInterface, err := client.NRClient.Nrqlmacros.GetEntity(macroID)
		if err != nil {
			log.Fatal(err)
		}

		if entityInterface == nil {
			log.Fatalf("no macro found with ID: %s", macroID)
		}

		macroEntity, ok := (*entityInterface).(*nrqlmacros.EntityManagementNRQLMacroEntity)
		if !ok {
			log.Fatalf("entity with ID %s is not a NRQL Macro", macroID)
		}

		err = output.Print(macroEntity)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosList = &cobra.Command{
	Use:   "list",
	Short: "List NRQL Macros",
	Long: `List NRQL Macros

Lists all NRQL Macros in your organization.
`,
	Example: "newrelic nrql macros list",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Nrqlmacros.GetEntitySearch("", "type = 'NRQL_MACRO'")
		if err != nil {
			log.Fatal(err)
		}

		if result == nil || len(result.Entities) == 0 {
			log.Info("No NRQL Macros found in your organization.")
			return
		}

		macros := make([]*nrqlmacros.EntityManagementNRQLMacroEntity, 0, len(result.Entities))
		for _, e := range result.Entities {
			m, ok := e.(*nrqlmacros.EntityManagementNRQLMacroEntity)
			if ok {
				macros = append(macros, m)
			}
		}

		err = output.Print(macros)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update a NRQL Macro",
	Long: `Update a NRQL Macro

Updates the expression or parameter count of an existing NRQL Macro.
Note: the macro name and scope cannot be changed after creation.

WARNING: Expression changes propagate immediately to all queries that
reference this macro. Ensure the updated expression is valid before saving.
`,
	Example: "newrelic nrql macros update --id <macroID> --expression \"environment = 'staging'\"",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		updateInput := nrqlmacros.EntityManagementNRQLMacroEntityUpdateInput{}

		if cmd.Flags().Changed("expression") {
			updateInput.Expression = macroExpression
		}

		if cmd.Flags().Changed("num-parameters") {
			updateInput.NumParameters = macroNumParams
		}

		if !cmd.Flags().Changed("expression") && !cmd.Flags().Changed("num-parameters") {
			log.Fatal("at least one of --expression or --num-parameters must be specified")
		}

		result, err := client.NRClient.Nrqlmacros.EntityManagementUpdateNRQLMacro(macroID, updateInput)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Print(result)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a NRQL Macro",
	Long: `Delete a NRQL Macro

Permanently deletes a NRQL Macro by its entity ID. Any queries that
reference this macro via $macroName will stop expanding the macro after deletion.
`,
	Example: "newrelic nrql macros delete --id <macroID>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := client.NRClient.Nrqlmacros.EntityManagementDelete(macroID)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("NRQL Macro with ID %s deleted successfully.", macroID)
	},
}

func init() {
	Command.AddCommand(cmdMacros)

	// create flags
	cmdMacrosCreate.Flags().StringVarP(&macroName, "name", "n", "", "Name of the NRQL Macro (required, immutable)")
	cmdMacrosCreate.Flags().StringVarP(&macroExpression, "expression", "e", "", "NRQL expression fragment for the macro (required)")
	cmdMacrosCreate.Flags().StringVar(&macroOrgID, "org-id", "", "Organization ID (auto-resolved from API if not provided)")
	cmdMacrosCreate.Flags().IntVar(&macroNumParams, "num-parameters", 0, "Number of positional parameters ($1, $2, ...) used in the expression")
	_ = cmdMacrosCreate.MarkFlagRequired("name")
	_ = cmdMacrosCreate.MarkFlagRequired("expression")

	// get flags
	cmdMacrosGet.Flags().StringVar(&macroID, "id", "", "ID of the NRQL Macro to retrieve (required)")
	_ = cmdMacrosGet.MarkFlagRequired("id")

	// list flags
	cmdMacrosList.Flags().StringVar(&macroOrgID, "org-id", "", "Organization ID (informational; macros are auto-filtered by your authenticated org)")

	// update flags
	cmdMacrosUpdate.Flags().StringVar(&macroID, "id", "", "ID of the NRQL Macro to update (required)")
	cmdMacrosUpdate.Flags().StringVarP(&macroExpression, "expression", "e", "", "Updated NRQL expression fragment")
	cmdMacrosUpdate.Flags().IntVar(&macroNumParams, "num-parameters", 0, "Updated number of positional parameters")
	_ = cmdMacrosUpdate.MarkFlagRequired("id")

	// delete flags
	cmdMacrosDelete.Flags().StringVar(&macroID, "id", "", "ID of the NRQL Macro to delete (required)")
	_ = cmdMacrosDelete.MarkFlagRequired("id")

	cmdMacros.AddCommand(cmdMacrosCreate)
	cmdMacros.AddCommand(cmdMacrosGet)
	cmdMacros.AddCommand(cmdMacrosList)
	cmdMacros.AddCommand(cmdMacrosUpdate)
	cmdMacros.AddCommand(cmdMacrosDelete)
}
