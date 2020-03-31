package nerdstorage

import (
	"fmt"
	"strings"

	"github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

var cmdCollection = &cobra.Command{
	Use:     "collection",
	Short:   "Read, write, and delete NerdStorage collections.",
	Example: "newrelic nerdstorage collection --help",
	Long:    `Read write, and delete NerdStorage collections`,
}

var cmdCollectionGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a NerdStorage collection.",
	Long: `Retrieve a NerdStorage collection

Retrieve a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `newrelic nerdstorage collection get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var resp []interface{}
			var err error

			input := nerdstorage.GetCollectionInput{
				PackageID:  packageID,
				Collection: collection,
			}

			switch strings.ToLower(scope) {
			case "account":
				if accountID == 0 {
					log.Fatal("account ID is required when using account scope")
				}

				resp, err = nrClient.NerdStorage.GetCollectionWithAccountScope(accountID, input)
			case "entity":
				if entityGUID == "" {
					log.Fatal("entity GUID is required when using entity scope")
				}

				resp, err = nrClient.NerdStorage.GetCollectionWithEntityScope(entityGUID, input)
			case "user":
				resp, err = nrClient.NerdStorage.GetCollectionWithUserScope(input)
			default:
				log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
			}
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(resp)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
			fmt.Println(string(json))
		})
	},
}

var cmdCollectionDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a NerdStorage collection.",
	Long: `Delete a NerdStorage collection

Delete a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `newrelic nerdstorage collection delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var err error

			input := nerdstorage.DeleteCollectionInput{
				PackageID:  packageID,
				Collection: collection,
			}

			switch strings.ToLower(scope) {
			case "account":
				if accountID == 0 {
					log.Fatal("account ID is required when using account scope")
				}

				_, err = nrClient.NerdStorage.DeleteCollectionWithAccountScope(accountID, input)
			case "entity":
				if entityGUID == "" {
					log.Fatal("entity GUID is required when using entity scope")
				}

				_, err = nrClient.NerdStorage.DeleteCollectionWithEntityScope(entityGUID, input)
			case "user":
				_, err = nrClient.NerdStorage.DeleteCollectionWithUserScope(input)
			default:
				log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
			}
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
	},
}

func init() {
	Command.AddCommand(cmdCollection)

	cmdCollection.AddCommand(cmdCollectionGet)
	cmdCollectionGet.Flags().IntVar(&accountID, "accountId", 0, "the account ID")
	cmdCollectionGet.Flags().StringVar(&entityGUID, "entityGuid", "", "the entity GUID")
	cmdCollectionGet.Flags().StringVar(&packageID, "packageId", "", "the external package ID")
	cmdCollectionGet.Flags().StringVar(&collection, "collection", "", "the collection name to get the document from")
	cmdCollectionGet.Flags().StringVar(&scope, "scope", "USER", "the scope to get the document from")

	err := cmdCollectionGet.MarkFlagRequired("packageId")
	if err != nil {
		log.Error(err)
	}

	err = cmdCollectionGet.MarkFlagRequired("scope")
	if err != nil {
		log.Error(err)
	}

	err = cmdCollectionGet.MarkFlagRequired("collection")
	if err != nil {
		log.Error(err)
	}

	cmdCollection.AddCommand(cmdCollectionDelete)
	cmdCollectionDelete.Flags().IntVar(&accountID, "accountId", 0, "the account ID")
	cmdCollectionDelete.Flags().StringVar(&entityGUID, "entityGuid", "", "the entity GUID")
	cmdCollectionDelete.Flags().StringVar(&packageID, "packageId", "", "the external package ID")
	cmdCollectionDelete.Flags().StringVar(&collection, "collection", "", "the collection name to delete the document from")
	cmdCollectionDelete.Flags().StringVar(&scope, "scope", "USER", "the scope to delete the document from")

	err = cmdCollectionDelete.MarkFlagRequired("packageId")
	if err != nil {
		log.Error(err)
	}

	err = cmdCollectionDelete.MarkFlagRequired("scope")
	if err != nil {
		log.Error(err)
	}

	err = cmdCollectionDelete.MarkFlagRequired("collection")
	if err != nil {
		log.Error(err)
	}
}
