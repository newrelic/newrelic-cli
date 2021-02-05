package decode

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	paramElement  string
	searchElement string
)

// Command represents the decode command.
var cmdDecode = &cobra.Command{
	Use:     "url",
	Short:   "Decodes NR1 URL Strings ",
	Example: `newrelic decode url -p="pane" -s="entityId" https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1`,

	Run: func(cmd *cobra.Command, args []string) {
		relicString := strings.Join(args, "")
		u, _ := url.Parse(relicString)
		relicString = u.Query().Get(paramElement)
		decoded, err := base64.StdEncoding.DecodeString(relicString)

		utils.LogIfFatal(err)

		decodedURL := string(decoded)
		entityID := strings.Contains(decodedURL, searchElement)

		if !entityID {
			fmt.Printf("%s not found in %s", searchElement, decodedURL)
			return
		}

		entityIDString := gjson.Get(decodedURL, searchElement)
		decodedEntityID := entityIDString.String()

		var regex = regexp.MustCompile(`[a-zA-Z0-9\+]*={0,3}`)
		//regex is used to determine correct base64 length
		if len(regex.FindStringIndex(decodedEntityID)) > 0 {
			for (len(decodedEntityID) % 4) != 0 {
				decodedEntityID += "="
			}
		}
		entityVal := string(decoded)

		decoded, err = base64.StdEncoding.DecodeString(decodedEntityID)
		if err != nil {
			entityIDString := gjson.Get(entityVal, searchElement)
			fmt.Println(entityIDString)
		}

		decodedEntityID = string(decoded)
		fmt.Println(decodedEntityID)

	},
}

func init() {
	Command.AddCommand(cmdDecode)
	cmdDecode.Flags().StringVarP(&paramElement, "param", "p", "", "the query parameter you want to decode")
	utils.LogIfError(cmdDecode.MarkFlagRequired("param"))

	cmdDecode.Flags().StringVarP(&searchElement, "search", "s", "", "the search key you want returned")
	utils.LogIfError(cmdDecode.MarkFlagRequired("search"))
}
