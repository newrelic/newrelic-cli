package extensions

import (
	"fmt"
	"os"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "extensions",
	Short: "Testing extensions",
}

var cmdHello = &cobra.Command{
	Use:     "hello",
	Short:   "Hello?",
	Long:    `Say hello to my little friend`,
	Example: "newrelic hello",
	Run: func(cmd *cobra.Command, args []string) {
		Do(cmd, args)
	},
}

func init() {
	Command.AddCommand(cmdHello)
}

func Do(cmd *cobra.Command, args []string) {
	go func() {
		fmt.Print("\nStarting server... \n")
		serve(cmd, args)
	}()

	manifest := Manifest{
		Command: "sleep",
	}

	proc, err := New(&manifest,
		WithTimeout(time.Duration(1000)),
		WithArgs("11"),
	)
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	proc.Stdout(os.Stdout)
	proc.Stdin(os.Stdin)
	proc.Stderr(os.Stderr)

	err = proc.Start()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	<-proc.DoneChan

	procErr := proc.Err()

	typeofType := reflect.TypeOf(procErr)

	switch procErr.(type) {
	case *ErrorDeadlineExceeded:
		log.Fatalf("Error: DeadlineExceeded: %+v", procErr)
	case *ErrorExit:
		log.Fatalf("Error: ExitError: %+v", procErr)
	default:
		log.Print("You are a good developer.")
	}

	fmt.Print("\nYou are a good developer.\n")
}
