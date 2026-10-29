package nrql

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqlmacros"
)

var (
	macroID             string
	macroName           string
	macroExpression     string
	macroOrgID          string
	macroParameterCount int
)

var cmdMacros = &cobra.Command{
	Use:   "macros",
	Short: "Manage NRQL Macros",
	Long: `Manage NRQL Macros in your New Relic organization.

A NRQL Macro is a named, reusable expression fragment. Wherever you write
$macroName in a NRQL query, New Relic substitutes the macro's expression
before running the query. Macros are available across every account in your
organization.`,
	Example: "newrelic nrql macros --help",
}

var cmdMacrosCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a NRQL Macro",
	Long: `Create a NRQL Macro in your organization.

The macro name is permanent — it cannot be changed after creation. Choose
a name that reflects the macro's purpose, because that name becomes part of
every query that references it ($macroName).

When the expression uses positional parameters ($1, $2, ...), you must also
pass --parameter-count with the exact number of distinct tokens.`,
	Example: `  # Simple filter macro
  newrelic nrql macros create --name filterProd --expression "environment = 'production'"

  # Parameterized macro
  newrelic nrql macros create --name filterStatusRegion --expression "status = $1 AND region = $2" --parameter-count 2`,
	PreRun: client.RequireClient,
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
				Name:          macroName,
				Expression:    macroExpression,
				NumParameters: macroParameterCount,
				Scope: nrqlmacros.EntityManagementScopedReferenceInput{
					Type: nrqlmacros.EntityManagementEntityScopeTypes.ORGANIZATION,
					ID:   orgID,
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Print(result.Entity)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosGet = &cobra.Command{
	Use:   "get",
	Short: "Fetch a NRQL Macro by ID",
	Long: `Fetch a single NRQL Macro by its entity GUID.

The ID is the entity GUID returned when you create a macro or list all macros.`,
	Example: "  newrelic nrql macros get --id <macroID>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		entityInterface, err := client.NRClient.Nrqlmacros.GetEntity(macroID)
		if err != nil {
			log.Fatal(err)
		}

		if entityInterface == nil {
			log.Fatalf("no macro found with ID %s", macroID)
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
	Short: "List all NRQL Macros in your organization",
	Long: `List all NRQL Macros that belong to your organization.

The list is scoped to the organization associated with your API key.`,
	Example: "  newrelic nrql macros list",
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
			if m, ok := e.(*nrqlmacros.EntityManagementNRQLMacroEntity); ok {
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
	Long: `Update the expression or parameter count of an existing NRQL Macro.

The macro name cannot be changed after creation. Provide at least one of
--expression or --parameter-count.

WARNING: Expression changes take effect immediately. Every query, dashboard
widget, and alert condition that references this macro will use the new
expression on its next execution.`,
	Example: `  # Update the expression
  newrelic nrql macros update --id <macroID> --expression "environment = 'staging'"

  # Clear positional parameters
  newrelic nrql macros update --id <macroID> --expression "environment = 'prod'" --parameter-count 0`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("expression") && !cmd.Flags().Changed("parameter-count") {
			log.Fatal("at least one of --expression or --parameter-count must be specified")
		}

		updateInput := nrqlmacros.EntityManagementNRQLMacroEntityUpdateInput{}

		if cmd.Flags().Changed("expression") {
			if macroExpression == "" {
				log.Fatal("--expression must not be empty")
			}
			updateInput.Expression = macroExpression
		}

		if cmd.Flags().Changed("parameter-count") {
			n := macroParameterCount
			updateInput.NumParameters = &n
		}

		result, err := client.NRClient.Nrqlmacros.EntityManagementUpdateNRQLMacro(macroID, updateInput)
		if err != nil {
			log.Fatal(err)
		}

		err = output.Print(result.Entity)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdMacrosDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a NRQL Macro",
	Long: `Permanently delete a NRQL Macro.

After deletion, any query that references this macro with $macroName will
no longer expand it. This action cannot be undone.`,
	Example: "  newrelic nrql macros delete --id <macroID>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := client.NRClient.Nrqlmacros.EntityManagementDelete(macroID)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("NRQL Macro %s deleted.", macroID)
	},
}

func init() {
	Command.AddCommand(cmdMacros)

	// create
	cmdMacrosCreate.Flags().StringVarP(&macroName, "name", "n", "", "Unique name for the macro — used as $name in NRQL queries (required, immutable)")
	cmdMacrosCreate.Flags().StringVarP(&macroExpression, "expression", "e", "", "NRQL expression fragment the macro expands to (required)")
	cmdMacrosCreate.Flags().StringVar(&macroOrgID, "org-id", "", "Organization GUID (resolved automatically from your API key if omitted)")
	cmdMacrosCreate.Flags().IntVar(&macroParameterCount, "parameter-count", 0, "Number of positional parameters ($1, $2, ...) in the expression; required when the expression contains $N tokens")
	_ = cmdMacrosCreate.MarkFlagRequired("name")
	_ = cmdMacrosCreate.MarkFlagRequired("expression")

	// get
	cmdMacrosGet.Flags().StringVar(&macroID, "id", "", "Entity GUID of the macro to fetch (required)")
	_ = cmdMacrosGet.MarkFlagRequired("id")

	// list — no flags; org is always inferred from the API key
	// (no --org-id: the entity search always returns macros for your authenticated organization)

	// update
	cmdMacrosUpdate.Flags().StringVar(&macroID, "id", "", "Entity GUID of the macro to update (required)")
	cmdMacrosUpdate.Flags().StringVarP(&macroExpression, "expression", "e", "", "New NRQL expression fragment")
	cmdMacrosUpdate.Flags().IntVar(&macroParameterCount, "parameter-count", 0, "New positional parameter count; pass 0 to remove all parameters")
	_ = cmdMacrosUpdate.MarkFlagRequired("id")

	// delete
	cmdMacrosDelete.Flags().StringVar(&macroID, "id", "", "Entity GUID of the macro to delete (required)")
	_ = cmdMacrosDelete.MarkFlagRequired("id")

	cmdMacros.AddCommand(cmdMacrosCreate)
	cmdMacros.AddCommand(cmdMacrosGet)
	cmdMacros.AddCommand(cmdMacrosList)
	cmdMacros.AddCommand(cmdMacrosUpdate)
	cmdMacros.AddCommand(cmdMacrosDelete)
}
