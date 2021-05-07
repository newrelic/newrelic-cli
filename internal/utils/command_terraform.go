package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	label       string
	file        string
	outFile     string
	snakeCaseRE = regexp.MustCompile("^[a-z]+(_[a-z]+)*$")
)

var cmdTerraform = &cobra.Command{
	Use:   "terraform",
	Short: "Tools for working with Terraform",
	Long: `Tools for working with Terraform

The terraform commands can be used for generating Terraform HCL for simple observability
as code use cases.
`,
	Example: `cat terraform.json | newrelic utils terraform dashboard --label my_dashboard_resource`,
}

var cmdTerraformDashboard = &cobra.Command{
	Use:   "dashboard",
	Short: "Generate HCL for the newrelic_one_dashboard resource",
	Long: `Generate HCL for the newrelic_one_dashboard resource

This command generates HCL configuration for newrelic_one_dashboard resources from
exported JSON documents.  For more detail on exporting dashboards to JSON, see
https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/manage-your-dashboard/#dash-json

Input can be sourced from STDIN per the provided example, or from a file using the --file option.
Output will be sent to STDOUT by default but can be redirected to a file with the --out option.
`,
	Example: `cat terraform.json | newrelic utils terraform dashboard --label my_dashboard_resource`,
	Args: func(cmd *cobra.Command, args []string) error {
		if ok := snakeCaseRE.MatchString(label); !ok {
			return fmt.Errorf("resource label must be formatted with snake case: %s", label)
		}

		if file != "" {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", file)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var input []byte
		var err error
		if file != "" {
			input, err = ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			input, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
		}

		hcl := generateDashboardHCL(label, input)

		if outFile != "" {
			if err := ioutil.WriteFile(outFile, []byte(hcl), 0644); err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		} else {
			fmt.Print(hcl)
		}
	},
}

func init() {
	Command.AddCommand(cmdTerraform)

	cmdTerraform.AddCommand(cmdTerraformDashboard)
	cmdTerraformDashboard.Flags().StringVarP(&label, "label", "l", "", "the resource label to use when generating resource HCL")
	if err := cmdTerraformDashboard.MarkFlagRequired("label"); err != nil {
		log.Error(err)
	}

	cmdTerraformDashboard.Flags().StringVarP(&file, "file", "f", "", "a file that contains exported dashboard JSON")
	cmdTerraformDashboard.Flags().StringVarP(&outFile, "out", "o", "", "the file to send the generated HCL to")
}
