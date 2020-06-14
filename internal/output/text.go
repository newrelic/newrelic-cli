package output

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (o *Output) text(data interface{}) error {
	// Early quit on no data
	if data == nil {
		return nil
	}

	if o == nil {
		return errors.New("invalid output formatter")
	}

	// Let's see what they sent us
	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.String:
		fmt.Println(data)
	case reflect.Slice:
		return o.table(data)
	//case reflect.Struct:
	default:
		return fmt.Errorf("unable to format data type: %T", data)
	}

	return nil
}

func (o *Output) table(data interface{}) error {
	// Early quit on no data
	if data == nil {
		return nil
	}

	if o == nil {
		return errors.New("invalid output formatter")
	}

	tw := o.newTableWriter()

	// Let's see what they sent us
	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.Slice:
		// Create the header from the field names
		typ := reflect.TypeOf(data).Elem()

		//cols := first.NumField()
		cols := typ.NumField()
		header := make([]interface{}, cols)
		colConfig := make([]table.ColumnConfig, cols)

		for i := 0; i < cols; i++ {
			header[i] = typ.Field(i).Name
			colConfig[i].Name = typ.Field(i).Name
			colConfig[i].WidthMin = len(typ.Field(i).Name)
			colConfig[i].WidthMax = o.terminalWidth * 3 / 4
			colConfig[i].WidthMaxEnforcer = text.WrapSoft
		}
		tw.SetColumnConfigs(colConfig)
		tw.AppendHeader(table.Row(header))

		// Add all the rows
		for i := 0; i < v.Len(); i++ {
			row := make([]interface{}, v.Index(i).NumField())
			for f := 0; f < v.Index(i).NumField(); f++ {
				row[f] = v.Index(i).Field(f).Interface()
			}
			tw.AppendRow(table.Row(row))
		}

	//case reflect.Struct:
	default:
		return fmt.Errorf("unable to format data type: %T", data)
	}

	tw.Render()

	return nil
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

	return t
}
