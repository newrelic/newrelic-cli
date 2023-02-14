package utils

// Borrowed from gojq - cli package
//
// The ability to do yaml processing is internal to gojq, so here
// we've copied in the nesessary bits to get the bare functionality.
// With this implementation, this subcommand will function much like
// yq: https://github.com/mikefarah/yq
//
// ref: https://github.com/itchyny/gojq/blob/main/cli/yaml.go

import (
	"fmt"
	"os"

	"github.com/itchyny/gojq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	yq "github.com/newrelic/newrelic-cli/internal/utils/yq"
)

var (
	outputIndent = 2
)

var cmdYq = &cobra.Command{
	Use:   "yq",
	Short: "Parse yaml strings",
	Long: `Parse yaml strings

The yq subcommand makes use of gojq (https://github.com/itchyny/gojq) to provide
yaml parsing capabilities.
`,
	Example: `echo '"foo": 128' | newrelic utils yq '.foo'`,
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

		iter := yq.NewYAMLInputIter(os.Stdin, "<stdin>")
		code, err := gojq.Compile(
			query,
			gojq.WithInputIter(iter),
		)
		if err != nil {
			log.Fatalln(err)
		}

		err = process(iter, code)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

// cli.process
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L337
func process(iter yq.InputIter, code *gojq.Code) error {
	var err error
	for {
		v, ok := iter.Next()
		if !ok {
			return err
		}
		if er, ok := v.(error); ok {
			printError(er)
			err = &yq.EmptyError{Err: er}
			// log.Fatalln(er) //todo ?
			continue
		}
		// TODO: if er := cli.printValues(code.Run(v, cli.argvalues...)); er != nil {
		if er := printValues(code.Run(v)); er != nil {
			printError(er)
			err = &yq.EmptyError{Err: er}
		}
	}
}

// cli.printValues
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L356
func printValues(iter gojq.Iter) error {
	m := createMarshaler()
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}

		if err := m.Marshal(v, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

// cli.printError
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L422
func printError(err error) {
	if er, ok := err.(interface{ IsEmptyError() bool }); !ok || !er.IsEmptyError() {
		if er, ok := err.(interface{ IsHaltError() bool }); !ok || !er.IsHaltError() {
			fmt.Fprintf(os.Stderr, "%s: %s\n", "cmdYq", err)
		} else if er, ok := err.(gojq.ValueError); ok {
			v := er.Value()
			if str, ok := v.(string); ok {
				os.Stderr.Write([]byte(str))
			} else {
				bs, _ := gojq.Marshal(v)
				os.Stderr.Write(bs)
				os.Stderr.Write([]byte{'\n'})
			}
		}
	}
}

// cli.createMarshaler
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L392
func createMarshaler() yq.Marshaler {
	return yq.YamlFormatter(&outputIndent)
}

func init() {
	Command.AddCommand(cmdYq)
}
