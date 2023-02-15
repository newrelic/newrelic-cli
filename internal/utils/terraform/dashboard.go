package terraform

import (
	"encoding/json"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
)

var (
	dashboardResourceName = "newrelic_one_dashboard"
	widgetTypes           = map[string]string{
		"viz.area":                "widget_area",
		"viz.bar":                 "widget_bar",
		"viz.billboard":           "widget_billboard",
		"viz.bullet":              "widget_bullet",
		"viz.funnel":              "widget_funnel",
		"viz.heatmap":             "widget_heatmap",
		"viz.histogram":           "widget_histogram",
		"viz.json":                "widget_json",
		"viz.line":                "widget_line",
		"viz.markdown":            "widget_markdown",
		"viz.pie":                 "widget_pie",
		"viz.table":               "widget_table",
		"viz.stacked-bar":         "widget_stacked_bar",
		"logger.log-table-widget": "widget_log_table",
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
	NullValues        DashboardWidgetNullValues  `json:"nullValues"`
	Units             DashboardWidgetUnits       `json:"units"`
	Colors            DashboardWidgetColors      `json:"colors"`
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

type DashboardWidgetNullValues struct {
	NullValue       string                              `json:"nullValue"`
	SeriesOverrides []DashboardWidgetNullValueOverrides `json:"seriesOverrides"`
}

type DashboardWidgetNullValueOverrides struct {
	NullValue  string `json:"nullValue"`
	SeriesName string `json:"seriesName"`
}
type DashboardWidgetUnits struct {
	Unit            string                         `json:"unit"`
	SeriesOverrides []DashboardWidgetUnitOverrides `json:"seriesOverrides"`
}

type DashboardWidgetUnitOverrides struct {
	Unit       string `json:"unit"`
	SeriesName string `json:"seriesName"`
}

type DashboardWidgetColors struct {
	Color           string                          `json:"color"`
	SeriesOverrides []DashboardWidgetColorOverrides `json:"seriesOverrides"`
}

type DashboardWidgetColorOverrides struct {
	Color      string `json:"color"`
	SeriesName string `json:"seriesName"`
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
						h.WriteBlock("facet", []string{}, func() {
							h.WriteStringAttribute("showOtherSeries", strconv.FormatBool(config.Facet.ShowOtherSeries))
						})
						h.WriteBlock("legend", []string{}, func() {
							h.WriteStringAttribute("enabled", strconv.FormatBool(config.Legend.Enabled))
						})
						h.WriteBlock("yAxisLeft", []string{}, func() {
							h.WriteStringAttribute("zero", strconv.FormatBool(config.YAxisLeft.Zero))
						})
						h.WriteBlock("nullValues", []string{}, func() {
							h.WriteStringAttribute("nullValue", config.NullValues.NullValue)
							for _, so := range config.NullValues.SeriesOverrides {
								h.WriteBlock("seriesOverrides", []string{}, func() {
									h.WriteStringAttribute("nullValue", so.NullValue)
									h.WriteStringAttribute("seriesName", so.SeriesName)
								})
							}
						})
						h.WriteBlock("units", []string{}, func() {
							h.WriteStringAttribute("unit", config.Units.Unit)
							for _, so := range config.Units.SeriesOverrides {
								h.WriteBlock("seriesOverrides", []string{}, func() {
									h.WriteStringAttribute("unit", so.Unit)
									h.WriteStringAttribute("seriesName", so.SeriesName)
								})
							}
						})
						h.WriteBlock("colors", []string{}, func() {
							h.WriteStringAttribute("color", config.Colors.Color)
							for _, so := range config.Colors.SeriesOverrides {
								h.WriteBlock("seriesOverrides", []string{}, func() {
									h.WriteStringAttribute("color", so.Color)
									h.WriteStringAttribute("seriesName", so.SeriesName)
								})
							}
						})

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
