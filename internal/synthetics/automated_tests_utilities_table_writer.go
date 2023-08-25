package synthetics

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Output struct {
	terminalWidth int
}

func printResultTable(tableData [][]string) {
	o := &Output{terminalWidth: 100}

	tw := o.newTableWriter()
	tw.Style().Name = "nr-syn-cli-table"

	// Add the header
	tw.AppendHeader(table.Row{"Status", "Monitor Name", "Monitor GUID", "isBlocking"})

	// Add the rows
	for _, row := range tableData {
		tw.AppendRow(stringSliceToRow(row))
	}

	// Render the table
	tw.Render()
}

func stringSliceToRow(slice []string) table.Row {
	row := make(table.Row, len(slice))
	for i, v := range slice {
		row[i] = v
	}
	return row
}

func (o *Output) newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(o.terminalWidth)

	t.SetStyle(table.StyleRounded)
	t.SetStyle(table.Style{
		Name: "nr-cli-table",
		//Box:  table.StyleBoxRounded,
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
	t.SetStyle(table.Style{
		Name: "nr-syn-cli-table",
		Box:  table.StyleBoxRounded,
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})
	return t
}
