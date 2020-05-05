package extensions

import (
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// CommandFlag represents a command flag.
type CommandFlag struct {
	Name      string   `yaml:"Name,omitempty"`
	Options   []string `yaml:"Options,omitempty"`
	Prompt    string   `yaml:"Prompt,omitempty"`
	Required  bool     `yaml:"Required,omitempty"`
	Shorthand string   `yaml:"Shorthand,omitempty"`
	Type      string   `yaml:"Type,omitempty"`
	Usage     string   `yaml:"Usage,omitempty"`
}

// CommandDefinition represents the definition of a CLI subcommand.
type CommandDefinition struct {
	Use         string         `yaml:"Use,omitempty"`
	Short       string         `yaml:"Short,omitempty"`
	Long        string         `yaml:"Long,omitempty"`
	Flags       []*CommandFlag `yaml:"Flags,omitempty"`
	Interactive bool           `yaml:"Interactive,omitempty"`
}

// Flag retrieves a specific CommandFlag by name.
func (cd *CommandDefinition) Flag(name string) *CommandFlag {
	for _, f := range cd.Flags {
		if f.Name == name {
			return f
		}
	}

	return nil
}

func initializeExtensions() {
	// Discover extensions by reading their manifests
	extensions, err := discoverExtensions()
	if err != nil {
		log.Fatal(err)
	}

	for _, def := range extensions {
		err = initializeExtension(def)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func discoverExtensions() ([]*ExtensionManifest, error) {
	manifests := []*ExtensionManifest{}

	files, err := ioutil.ReadDir(extensionLocation)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			file := extensionLocation + "/" + f.Name() + "/extension.yml"
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, err
			}

			var manifest *ExtensionManifest

			err = yaml.Unmarshal(dat, &manifest)
			if err != nil {
				return nil, err
			}

			manifest.Extension.Command = os.ExpandEnv(manifest.Extension.Command)
			manifests = append(manifests, manifest)
		}
	}

	return manifests, nil
}

func initializeExtension(p *ExtensionManifest) error {
	rootExtensionCmd := &cobra.Command{
		Use:   p.Extension.Name,
		Short: p.Extension.Short,
		Long:  p.Extension.Long,
	}

	// Add the extension's root command
	ExtensionRootCommands = append(ExtensionRootCommands, rootExtensionCmd)

	// Add subcommands to the current extension's root command
	subCmds, err := buildCobraCommands(p, p.Commands)
	if err != nil {
		return err
	}

	for _, c := range subCmds {
		rootExtensionCmd.AddCommand(c)
	}

	return nil
}

func buildCobraCommands(p *ExtensionManifest, commands []*CommandDefinition) ([]*cobra.Command, error) {
	cmds := []*cobra.Command{}
	for _, c := range commands {
		cmd, err := createCobraCommand(p, c)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func createCobraCommand(p *ExtensionManifest, cd *CommandDefinition) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   cd.Use,
		Short: cd.Short,
		Long:  cd.Long,
		Run: func(cmd *cobra.Command, args []string) {
			runExtension(p, cmd, args)
		},
	}

	for _, f := range cd.Flags {
		switch f.Type {
		case "bool":
			cmd.PersistentFlags().BoolP(f.Name, f.Shorthand, false, f.Usage)
		case "string":
			cmd.PersistentFlags().StringP(f.Name, f.Shorthand, "", f.Usage)
		}

		if f.Required {
			err := cmd.MarkPersistentFlagRequired(f.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return cmd, nil
}

func promptForValue(f *pflag.Flag, def *CommandFlag) error {
	var result string
	var err error
	if len(def.Options) > 0 {
		result, err = promptForStringOption(def.Prompt, def.Options)
	} else {
		result, err = promptForString(def.Prompt)
	}

	if err != nil {
		return err
	}

	err = f.Value.Set(result)
	if err != nil {
		return err
	}

	return nil
}

func promptForStringOption(prompt string, options []string) (string, error) {
	s := promptui.Select{
		Label: prompt,
		Items: options,
	}

	_, result, err := s.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func promptForString(prompt string) (string, error) {
	p := promptui.Prompt{
		Label: prompt,
	}

	result, err := p.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func runExtension(m *ExtensionManifest, cmd *cobra.Command, args []string) {
	go func() {
		serve(cmd, args)
	}()

	proc, err := New(m)

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

	switch procErr.(type) {
	case *ErrorDeadlineExceeded:
		log.Fatalf("Error: DeadlineExceeded: %+v", procErr)
	case *ErrorExit:
		log.Fatalf("Error: ExitError: %+v", procErr)
	}
}
