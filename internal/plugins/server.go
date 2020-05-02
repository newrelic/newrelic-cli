package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandableResponse struct {
	Cmd   string
	Args  []string
	Flags map[string]string
}

func onCommandRequest(command *cobra.Command, args []string) string {
	fmt.Printf("\n onCommandRequest - command: %+v", command.Name())

	flags := map[string]string{}
	command.Flags().VisitAll(func(f *pflag.Flag) {
		// Need more checks here, but making moves
		switch f.Value.Type() {
		case "string":
			fmt.Printf("\n onCommandRequest - flag: %+v=%+v", f.Name, f.Value.String())
			flags[f.Name] = f.Value.String()
		}
	})

	fmt.Printf("\n onCommandRequest - flags: %+v", flags)

	returnObj := CommandableResponse{
		Cmd:   command.Name(),
		Args:  args,
		Flags: flags,
	}

	resp, _ := json.Marshal(returnObj)

	fmt.Printf("\n onCommandRequest - resp: %+v \n", string(resp))

	return string(resp)
}

func serve(command *cobra.Command, args []string) {
	fmt.Printf("\n\n serve - command: %s \n\n", command.Name())

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("\n\n handle request - command: %s \n\n", command.Name())

		result := onCommandRequest(command, args)

		fmt.Fprintf(w, "%s", result)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
