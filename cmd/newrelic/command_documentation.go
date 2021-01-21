package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docOutputDir string
	docFormat    string
)

var cmdDocumentation = &cobra.Command{
	Use:   "documentation",
	Short: "Generate CLI documentation",
	Long: `Generate CLI documentation

newrelic documentation --outputDir <my directory> --type (markdown|manpage)

`,
	Example: "newrelic documentation --outputDir /tmp",
	Run: func(cmd *cobra.Command, args []string) {
		if docOutputDir == "" {
			if err := cmd.Help(); err != nil {
				log.Warn(err)
			}
			log.Fatal("--outputDir <my directory> is required")
		}

		switch docFormat {
		case "markdown":
			err := doc.GenMarkdownTree(Command, docOutputDir)
			if err != nil {
				log.Error(err)
			}
		case "manpage":
			header := &doc.GenManHeader{
				Title:   "newrelic",
				Section: "3",
			}
			err := doc.GenManTree(Command, header, docOutputDir)
			if err != nil {
				log.Error(err)
			}
		default:
			if err := cmd.Help(); err != nil {
				log.Error(err)
			}
			log.Error("--type must be one of [markdown, manpage]")
		}
	},
}

func init() {
	Command.AddCommand(cmdDocumentation)

	cmdDocumentation.Flags().StringVarP(&docOutputDir, "outputDir", "o", "", "Output directory for generated documentation")
	cmdDocumentation.Flags().StringVar(&docFormat, "format", "markdown", "Documentation format [markdown, manpage] default 'markdown'")
	if err := cmdDocumentation.MarkFlagRequired("outputDir"); err != nil {
		log.Error(err)
	}
}
