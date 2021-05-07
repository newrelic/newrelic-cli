package terraform

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var (
	widgetTypes = map[string]string{
		"viz.area":      "widget_area",
		"viz.bar":       "widget_bar",
		"viz.billboard": "widget_billboard",
		"viz.bullet":    "widget_bullet",
		"viz.funnel":    "widget_funnel",
		"viz.heatmap":   "widget_heatmap",
		"viz.histogram": "widget_histogram",
		"viz.json":      "widget_json",
		"viz.line":      "widget_line",
		"viz.markdown":  "widget_markdown",
		"viz.pie":       "widget_pie",
		"viz.table":     "widget_table",
	}
)

type Dashboard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions string `json:"permissions"`

	Pages []DashboardPage `json:"pages"`
}

type DashboardPage struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Widgets []DashboardWidget `json:"widgets"`
}

type DashboardWidget struct {
	Title            string                          `json:"title"`
	Visualization    DashboardWidgetVisualization    `json:"visualization"`
	Layout           DashboardWidgetLayout           `json:"layout"`
	RawConfiguration DashboardWidgetRawConfiguration `json:"rawConfiguration"`
}

type DashboardWidgetVisualization struct {
	ID string `json:"id"`
}

type DashboardWidgetLayout struct {
	Column int `json:"column"`
	Row    int `json:"row"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

type DashboardWidgetRawConfiguration struct {
	DataFormatters    []string                   `json:"dataFormatters"`
	NRQLQueries       []DashboardWidgetNRQLQuery `json:"nrqlQueries"`
	LinkedEntityGUIDs []string                   `json:"linkedEntityGuids"`
	Text              string                     `json:"text"`
	Facet             DashboardWidgetFacet       `json:"facet"`
	Legend            DashboardWidgetLegend      `json:"legend"`
	YAxisLeft         DashboardWidgetYAxisLeft   `json:"yAxisLeft"`
}

type DashboardWidgetFacet struct {
	ShowOtherSeries bool `json:"showOtherSeries"`
}

type DashboardWidgetNRQLQuery struct {
	AccountID int    `json:"accountId"`
	Query     string `json:"query"`
}

type DashboardWidgetLegend struct {
	Enabled bool `json:"enabled"`
}

type DashboardWidgetYAxisLeft struct {
	Zero bool `json:"zero"`
}

func GenerateDashboardHCL(resourceLabel string, input []byte) (string, error) {
	var d Dashboard
	if err := json.Unmarshal(input, &d); err != nil {
		log.Fatal(err)
	}

	resourceName := "newrelic_one_dashboard"

	h := NewHCLGen(2)
	h.WriteBlock("resource", []string{resourceName, resourceLabel}, func() error {
		h.WriteStringAttribute("name", d.Name)
		h.WriteStringAttributeIfNotEmpty("description", d.Description)
		h.WriteStringAttributeIfNotEmpty("permissions", strings.ToLower(d.Permissions))

		for _, p := range d.Pages {
			h.WriteBlock("page", []string{}, func() error {
				h.WriteStringAttribute("name", p.Name)
				h.WriteStringAttributeIfNotEmpty("description", p.Description)

				for _, w := range p.Widgets {
					if widgetTypes[w.Visualization.ID] == "" {
						return fmt.Errorf("unrecognized widget type \"%s\"\n", w.Visualization.ID)
					}

					h.WriteBlock(widgetTypes[w.Visualization.ID], []string{}, func() error {
						h.WriteStringAttribute("title", w.Title)
						h.WriteIntAttribute("row", w.Layout.Row)
						h.WriteIntAttribute("column", w.Layout.Column)
						h.WriteIntAttribute("height", w.Layout.Height)
						h.WriteIntAttribute("width", w.Layout.Width)
						h.WriteStringSliceAttributeIfNotEmpty("linked_entity_guids", w.RawConfiguration.LinkedEntityGUIDs)
						h.WriteMultilineStringAttribute("text", w.RawConfiguration.Text)

						for _, q := range w.RawConfiguration.NRQLQueries {
							h.WriteBlock("nrql_query", []string{}, func() error {
								h.WriteIntAttributeIfNotZero("account_id", q.AccountID)
								h.WriteMultilineStringAttribute("query", q.Query)

								return nil
							})
						}

						return nil
					})
				}

				return nil
			})
		}

		return nil
	})

	return h.String(), nil
}
