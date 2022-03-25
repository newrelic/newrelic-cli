package terraform

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
)

var (
	dashboardResourceName = "newrelic_one_dashboard"
	widgetTypes           = map[string]string{
		"viz.area":        "widget_area",
		"viz.bar":         "widget_bar",
		"viz.billboard":   "widget_billboard",
		"viz.bullet":      "widget_bullet",
		"viz.funnel":      "widget_funnel",
		"viz.heatmap":     "widget_heatmap",
		"viz.histogram":   "widget_histogram",
		"viz.json":        "widget_json",
		"viz.line":        "widget_line",
		"viz.markdown":    "widget_markdown",
		"viz.pie":         "widget_pie",
		"viz.table":       "widget_table",
		"viz.stacked-bar": "widget_stacked_bar",
	}
)

type DashboardWidgetRawConfiguration struct {
	DataFormatters    []DataFormatter            `json:"dataFormatters"`
	NRQLQueries       []DashboardWidgetNRQLQuery `json:"nrqlQueries"`
	LinkedEntityGUIDs []string                   `json:"linkedEntityGuids"`
	Text              string                     `json:"text"`
	Facet             DashboardWidgetFacet       `json:"facet"`
	Legend            DashboardWidgetLegend      `json:"legend"`
	YAxisLeft         DashboardWidgetYAxisLeft   `json:"yAxisLeft"`
}

type DataFormatter struct {
	Name      string      `json:"name"`
	Precision interface{} `json:"precision"`
	Type      string      `json:"type"`
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

func GenerateDashboardHCL(resourceLabel string, shiftWidth int, input []byte) (string, error) {
	var d dashboards.DashboardInput
	if err := json.Unmarshal(input, &d); err != nil {
		log.Fatal(err)
	}

	h := NewHCLGen(shiftWidth)
	h.WriteBlock("resource", []string{dashboardResourceName, resourceLabel}, func() {
		h.WriteStringAttribute("name", d.Name)
		h.WriteStringAttributeIfNotEmpty("description", d.Description)
		h.WriteStringAttributeIfNotEmpty("permissions", strings.ToLower(string(d.Permissions)))

		for _, p := range d.Pages {
			h.WriteBlock("page", []string{}, func() {
				h.WriteStringAttribute("name", p.Name)
				h.WriteStringAttributeIfNotEmpty("description", p.Description)

				for _, w := range p.Widgets {
					requireValidVisualizationID(w.Visualization.ID)

					h.WriteBlock(widgetTypes[w.Visualization.ID], []string{}, func() {
						h.WriteStringAttribute("title", w.Title)
						h.WriteIntAttribute("row", w.Layout.Row)
						h.WriteIntAttribute("column", w.Layout.Column)
						h.WriteIntAttribute("height", w.Layout.Height)
						h.WriteIntAttribute("width", w.Layout.Width)

						config := unmarshalDashboardWidgetRawConfiguration(w.Title, widgetTypes[w.Visualization.ID], w.RawConfiguration)

						h.WriteStringSliceAttributeIfNotEmpty("linked_entity_guids", config.LinkedEntityGUIDs)
						h.WriteMultilineStringAttributeIfNotEmpty("text", config.Text)

						for _, q := range config.NRQLQueries {
							h.WriteBlock("nrql_query", []string{}, func() {
								h.WriteIntAttributeIfNotZero("account_id", q.AccountID)
								h.WriteMultilineStringAttribute("query", q.Query)
							})
						}
					})
				}
			})
		}
	})

	return h.String(), nil
}

func unmarshalDashboardWidgetRawConfiguration(title string, widgetType string, b []byte) *DashboardWidgetRawConfiguration {
	var c DashboardWidgetRawConfiguration
	err := json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("failed unmarshaling rawConfiguration for widget \"%s\" of type \"%s\"", title, widgetType)
		panic(err)
	}

	return &c
}

func requireValidVisualizationID(id string) {
	if widgetTypes[id] == "" {
		log.Fatalf("unrecognized widget type \"%s\"", id)
	}
}
