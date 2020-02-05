package entities

import (
	"fmt"
	"log"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var (
	nrClient     *newrelic.NewRelic
	entityName   string
	entityGUID   string
	entityTags   []string
	entityValues []string
)

// SetClient is the API for passing along the New Relic client to this command
func SetClient(nr *newrelic.NewRelic) error {
	if nr == nil {
		return fmt.Errorf("client can not be nil")
	}

	nrClient = nr

	return nil
}

// Command represents the entities command
var Command = &cobra.Command{
	Use:   "entities",
	Short: "entities commands",
}

var entitiesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "entities search",
	Long: `Search for New Relic entities

The search command performs a search for New Relic eitities. Optionally, the
--name flag can be used to limit results to those that match the the given
name.
`,
	Example: "newrelic entities search -n test",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		params := entities.SearchEntitiesParams{
			Name: entityName,
		}

		entities, err := nrClient.Entities.SearchEntities(params)
		if err != nil {
			log.Fatal(err)
		}

		json, err := prettyjson.Marshal(entities)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	},
}

var entitiesDescribeTags = &cobra.Command{
	Use:   "describe-tags",
	Short: "describe-tags",
	Long: `Describe the tags for a given entity

The describe-tags command returns JSON output of the tags for the requested
entity.
`,
	Example: "newrelic entities describe-tags --guid <guid>",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		tags, err := nrClient.Entities.ListTags(entityGUID)
		if err != nil {
			log.Fatal(err)
		}

		json, err := prettyjson.Marshal(tags)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	},
}

var entitiesDeleteTags = &cobra.Command{
	Use:   "delete-tags",
	Short: "delete-tags",
	Long: `Delete the given tags for the given entity

The delete-tags command performs a delete operation on the entity for the
specified tag keys.
`,
	Example: "newrelic entities delete-tags --guid <guid> --tag tag1 --tag tag2 --tag tag3,tag4",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		err := nrClient.Entities.DeleteTags(entityGUID, entityTags)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var entitiesDeleteTagValues = &cobra.Command{
	Use:   "delete-tag-values",
	Short: "delete-tag-values",
	Long: `Delete the given tag/value pairs from the given entitiy

The delete-tag-values command performs a delete operation on the entity for the
specified tag:value pairs.
`,
	Example: "newrelic entities delete-tag-values --guid <guid> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		tagValues, err := assembleTagValues(entityValues)
		if err != nil {
			log.Fatal(err)
		}

		err = nrClient.Entities.DeleteTagValues(entityGUID, tagValues)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var entitiesCreateTags = &cobra.Command{
	Use:   "create-tags",
	Short: "create-tags",
	Long: `Create tag:value pairs for the given entitiy

The create-tags command adds tag:value pairs for the given entity.
`,
	Example: "newrelic entities create-tags --guid <guid> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		tags, err := assembleTags(entityTags)
		if err != nil {
			log.Fatal(err)
		}

		err = nrClient.Entities.AddTags(entityGUID, tags)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var entitiesReplaceTags = &cobra.Command{
	Use:   "replace-tags",
	Short: "replace-tags",
	Long: `Repaces tag:value pairs for the given entitiy

The replace-tags command replaces any existing tag:value pairs with those
provided for the given entity.
`,
	Example: "newrelic entities replace-tags --guid <guid> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		tags, err := assembleTags(entityTags)
		if err != nil {
			log.Fatal(err)
		}

		err = nrClient.Entities.ReplaceTags(entityGUID, tags)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func assembleTags(tags []string) ([]entities.Tag, error) {
	var t []entities.Tag

	tagBuilder := make(map[string][]string)

	for _, x := range tags {
		if !strings.Contains(x, ":") {
			return []entities.Tag{}, fmt.Errorf("Tags must be specified as colon seperated key:value pairs")
		}

		v := strings.SplitN(x, ":", 2)

		tagBuilder[v[0]] = append(tagBuilder[v[0]], v[1])
	}

	for k, v := range tagBuilder {
		tag := entities.Tag{
			Key:    k,
			Values: v,
		}

		t = append(t, tag)
	}

	return t, nil
}

func assembleTagValues(values []string) ([]entities.TagValue, error) {
	var tagValues []entities.TagValue

	for _, x := range values {
		if !strings.Contains(x, ":") {
			return []entities.TagValue{}, fmt.Errorf("Tag values must be specified as colon seperated key:value pairs")
		}

		v := strings.SplitN(x, ":", 2)

		tv := entities.TagValue{
			Key:   v[0],
			Value: v[1],
		}
		tagValues = append(tagValues, tv)
	}

	return tagValues, nil
}

func init() {
	Command.AddCommand(entitiesSearchCmd)
	entitiesSearchCmd.Flags().StringVarP(&entityName, "name", "n", "ENTITY_NAME", "search for results matching the given name")

	Command.AddCommand(entitiesDescribeTags)
	entitiesDescribeTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to describe")
	entitiesDescribeTags.MarkFlagRequired("guid")

	Command.AddCommand(entitiesDeleteTags)
	entitiesDeleteTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tags on")
	entitiesDeleteTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to delete from the entity")
	entitiesDeleteTags.MarkFlagRequired("guid")
	entitiesDeleteTags.MarkFlagRequired("tag")

	Command.AddCommand(entitiesDeleteTagValues)
	entitiesDeleteTagValues.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tag values on")
	entitiesDeleteTagValues.Flags().StringSliceVarP(&entityValues, "value", "v", []string{}, "keyy:value tags to delete from the entity")
	entitiesDeleteTagValues.MarkFlagRequired("guid")
	entitiesDeleteTagValues.MarkFlagRequired("value")

	Command.AddCommand(entitiesCreateTags)
	entitiesCreateTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tag values on")
	entitiesCreateTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to add to the entity")
	entitiesCreateTags.MarkFlagRequired("guid")
	entitiesCreateTags.MarkFlagRequired("tag")

	Command.AddCommand(entitiesReplaceTags)
	entitiesReplaceTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tag values on")
	entitiesReplaceTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to replace on the entity")
	entitiesReplaceTags.MarkFlagRequired("guid")
	entitiesReplaceTags.MarkFlagRequired("tag")
}
