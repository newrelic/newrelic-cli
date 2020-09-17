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
	case reflect.Slice, reflect.Struct:
		return o.renderAsTable(data)
	default:
		return fmt.Errorf("unable to format data type: %T", data)
	}

	return nil
}

func isPointer(element reflect.Value) bool {
	return element.Kind() == reflect.Ptr
}

func isMap(element reflect.Value) bool {
	return element.Kind() == reflect.Map
}

func (o *Output) renderAsTable(data interface{}) error {
	// Early quit on no data
	if data == nil {
		return nil
	}

	if o == nil {
		return errors.New("invalid output formatter")
	}

	tw := o.newTableWriter()

	fmt.Printf("\n data:  %v \n", data)

	// Let's see what they sent us
	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.Slice:
		// Return early if empty slice
		if v.Len() == 0 {
			return nil
		}

		arr := make([]reflect.Value, v.Len())

		for i := 0; i < v.Len(); i++ {
			el := v.Index(i)
			value := reflect.ValueOf(el)

			if isPointer(value) {
				arr[i] = reflect.Indirect(value).Elem()
			} else if isMap(el) {
				arr[i] = el
			} else {
				arr[i] = el.Elem()
			}
		}

		var cols int
		el := arr[0]
		if isMap(el) {
			for range el.MapKeys() {
				cols++
			}
		} else {
			cols = el.NumField()
		}

		fmt.Printf("\n el:  %v \n", el)
		fmt.Printf("\n cols:  %v \n", cols)

		// for _, e := range val.MapKeys() {
		// 	v := val.MapIndex(e)
		// 	switch t := v.Interface().(type) {
		// 	case int:
		// 		fmt.Println(e, t)
		// 	case string:
		// 		fmt.Println(e, t)
		// 	case bool:
		// 		fmt.Println(e, t)
		// 	default:
		// 		fmt.Println("not found")
		// 	}
		// }

		header := make([]interface{}, cols)
		colConfig := make([]table.ColumnConfig, cols)

		for i := 0; i < cols; i++ {
			var fld reflect.StructField
			if isMap(el) {
				fld = reflect.TypeOf(el).Field(i)
			} else {
				fld = el.Type().Field(i)
			}

			header[i] = fld.Name
			colConfig[i].Name = fld.Name
			colConfig[i].WidthMin = len(fld.Name)
			colConfig[i].WidthMax = o.terminalWidth * 3 / 4
			colConfig[i].WidthMaxEnforcer = text.WrapSoft
		}
		tw.SetColumnConfigs(colConfig)
		tw.AppendHeader(table.Row(header))

		// Add all the rows
		for i := 0; i < v.Len(); i++ {
			var numFields int
			element := v.Index(i)

			if isPointer(element) {
				numFields = reflect.Indirect(element).NumField()
			} else if isMap(element) {
				for range el.MapKeys() {
					numFields++
				}
			} else {
				numFields = element.NumField()
			}

			fmt.Printf("\n numFields:  %v \n", numFields)

			fmt.Print("\n\n **************************** \n")

			row := make([]interface{}, numFields)
			for f := 0; f < numFields; f++ {
				if isPointer(element) {
					elValue := reflect.Indirect(element)

					if isPointer(elValue.Field(f)) {
						row[f] = reflect.Indirect(elValue.Field(f)).Interface()
					} else {
						row[f] = elValue.Field(f).Interface()
					}
				} else {
					fmt.Printf("\n Row isValue:  %v - %v \n", element, element.Kind())

					if isMap(v.Index(i)) {
						row[f] = reflect.ValueOf(element).Field(i)
					} else {
						row[f] = v.Index(i).Field(f).Interface()
					}

					// row[f] = v.Index(i).Field(f).Interface()
				}

				// if isPointer(element) {
				// 	row[f] = reflect.Indirect(v.Index(i)).Field(f).Interface()
				// } else {
				// 	row[f] = v.Index(i).Field(f).Interface()
				// }

			}
			tw.AppendRow(table.Row(row))

			fmt.Print("\n **************************** \n\n")
		}

	// Single Struct becomes table view of Field | Value
	case reflect.Struct:
		typ := reflect.TypeOf(data)
		tw.AppendHeader(table.Row{"Field", "Value"})

		for f := 0; f < typ.NumField(); f++ {
			row := []interface{}{
				typ.Field(f).Name,
				v.Field(f).Interface(),
			}

			tw.AppendRow(table.Row(row))
		}

	default:
		return fmt.Errorf("unable to format data as table - type: %T", data)
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
