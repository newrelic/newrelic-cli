package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hokaccha/go-prettyjson"
	"github.com/itchyny/gojq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdJq = &cobra.Command{
	Use:   "jq",
	Short: "Parse json strings",
	Long: `Parse json strings

The jq subcommand makes use of gojq (https://github.com/itchyny/gojq) to provide
json parsing capabilities.
`,
	Example: `echo '{"foo": 128}' | newrelic utils jq '.foo'`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatalln("no filter string provided")
		}

		if !StdinExists() {
			log.Fatalln("no input found")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		query, err := gojq.Parse(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}

		var obj interface{}
		err = json.Unmarshal(bytes, &obj)
		if err != nil {
			log.Fatalln(err)
		}

		iter := query.Run(obj)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}

			if err, ok := v.(error); ok {
				log.Fatalln(err)
			}

			s, _ := prettyjson.Marshal(v)

			fmt.Println(string(s))
		}
	},
}

func init() {
	Command.AddCommand(cmdJq)
}
