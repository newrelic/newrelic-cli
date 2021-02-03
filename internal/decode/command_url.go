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
	paramElement string
	jsonElement  string
)

// Command represents the decode command.
var cmdDecode = &cobra.Command{
	Use:     "url",
	Short:   "Decodes NR1 URL Strings ",
	Example: `newrelic decode url -p="pane" -j="entityId" https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1`,

	Run: func(cmd *cobra.Command, args []string) {
		relicString := strings.Join(args, "")
		u, _ := url.Parse(relicString)
		relicString = u.Query().Get(paramElement)
		decoded, err := base64.StdEncoding.DecodeString(relicString)

		if err != nil {
			fmt.Println("decode error:", err)
		}

		decodedURL := string(decoded)
		fmt.Println(decodedURL)
		entityID := strings.Contains(decodedURL, jsonElement)

		if entityID {
			entityIDString := gjson.Get(decodedURL, jsonElement)
			decodedURL = entityIDString.String()

			var regex = regexp.MustCompile(`[a-zA-Z0-9\+]*={0,3}`)
			//regex is used to determine correct base64 length
			if len(regex.FindStringIndex(decodedURL)) > 0 {
				for (len(decodedURL) % 4) != 0 {
					decodedURL += "="
				}
			}
			entityVal := string(decoded)

			decoded, err := base64.StdEncoding.DecodeString(decodedURL)
			if err != nil {
				//fmt.Println("decode error:", err)
				entityIDString := gjson.Get(entityVal, jsonElement)
				fmt.Println(entityIDString)
			}

			decodedURL := string(decoded)
			fmt.Println(decodedURL)
		}
	},
}

func init() {
	Command.AddCommand(cmdDecode)
	cmdDecode.Flags().StringVarP(&paramElement, "param", "p", "", "The Param you want to decode")
	utils.LogIfError(cmdDecode.MarkFlagRequired("param"))

	cmdDecode.Flags().StringVarP(&jsonElement, "json", "j", "", "The Json element you want returned")
	utils.LogIfError(cmdDecode.MarkFlagRequired("json"))
}
