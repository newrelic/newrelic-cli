package config

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

// TableRenderer represents the presentation layer for human readable command output.
type TableRenderer struct{}

// NewTableRenderer returns a new instance of TableRenderer.
func NewTableRenderer() *TableRenderer {
	return &TableRenderer{}
}

// Get renders output for the get command.
func (tr *TableRenderer) Get(c *Config, key string) {
	tr.renderTable(c, key)
}

// List renders output for the list command.
func (tr *TableRenderer) List(c *Config) {
	tr.renderTable(c, "")
}

// Delete renders output for the delete command.
func (tr *TableRenderer) Delete(key string) {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s removed successfully\n", green("✔"), bold(key))
}

// Set renders output for the set command.
func (tr *TableRenderer) Set(key string, value string) {
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s set to %s\n", bold(key), cyan(value))
}

func (tr *TableRenderer) renderTable(c *Config, key string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value", "Origin"})
	t.AppendRows(tr.toRows(c, key))
	t.SetStyle(table.Style{
		Name: "nr-cli-table",
		Box: table.BoxStyle{
			MiddleHorizontal: "-",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:   "Value",
			Colors: text.Colors{text.FgHiCyan},
		},
		{
			Name:   "Origin",
			Colors: text.Colors{text.FgHiBlack},
		},
	})

	t.Render()

	bold := color.New(color.Bold).SprintFunc()
	fmt.Fprintf(os.Stderr, "\n\nRun %s for more info.\n", bold("\"newrelic config get --key KEY\""))
}

// toRows converts the config struct to go-pretty rows.
// Supplying a non-empty string value for name will limit the data
// to the field passed.
func (tr *TableRenderer) toRows(c *Config, name string) []table.Row {
	values := c.getAll(name)
	rows := []table.Row{}

	for _, v := range values {
		origin := "Default"
		if v.IsDefault() {
			origin = "User config"
		}

		row := table.Row{v.Name, v.Value, origin}
		rows = append(rows, row)
	}

	return rows
}
