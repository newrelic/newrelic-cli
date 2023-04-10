package terraform

import (
	"encoding/json"
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
	DataFormatters    []DataFormatter                `json:"dataFormatters"`
	NRQLQueries       []DashboardWidgetNRQLQuery     `json:"nrqlQueries"`
	LinkedEntityGUIDs []string                       `json:"linkedEntityGuids"`
	Text              string                         `json:"text"`
	Facet             DashboardWidgetFacet           `json:"facet,omitempty"`
	Legend            DashboardWidgetLegend          `json:"legend,omitempty"`
	YAxisLeft         DashboardWidgetYAxisLeft       `json:"yAxisLeft,omitempty"`
	NullValues        DashboardWidgetNullValues      `json:"nullValues,omitempty"`
	Units             DashboardWidgetUnits           `json:"units,omitempty"`
	Colors            DashboardWidgetColors          `json:"colors,omitempty"`
	PlatformOptions   DashboardWidgetPlatformOptions `json:"platformOptions,omitempty"`
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
	Enabled bool `json:"enabled,omitempty"`
}

type DashboardWidgetYAxisLeft struct {
	Max float64 `json:"max,omitempty"`
	Min float64 `json:"min,omitempty"`
}

type DashboardWidgetNullValues struct {
	NullValue       string                              `json:"nullValue,omitempty"`
	SeriesOverrides []DashboardWidgetNullValueOverrides `json:"seriesOverrides,omitempty"`
}

type DashboardWidgetNullValueOverrides struct {
	NullValue  string `json:"nullValue,omitempty"`
	SeriesName string `json:"seriesName,omitempty"`
}
type DashboardWidgetUnits struct {
	Unit            string                         `json:"unit,omitempty"`
	SeriesOverrides []DashboardWidgetUnitOverrides `json:"seriesOverrides,omitempty"`
}

type DashboardWidgetUnitOverrides struct {
	Unit       string `json:"unit,omitempty"`
	SeriesName string `json:"seriesName,omitempty"`
}

type DashboardWidgetColors struct {
	Color           string                          `json:"color,omitempty"`
	SeriesOverrides []DashboardWidgetColorOverrides `json:"seriesOverrides,omitempty"`
}

type DashboardWidgetColorOverrides struct {
	Color      string `json:"color,omitempty"`
	SeriesName string `json:"seriesName,omitempty"`
}
type DashboardWidgetPlatformOptions struct {
	IgnoreTimeRange bool `json:"ignoreTimeRange,omitempty"`
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
						h.WriteBooleanAttribute("facet_show_other_series", config.Facet.ShowOtherSeries)
						h.WriteBooleanAttribute("legend_enabled", config.Legend.Enabled)
						h.WriteBooleanAttribute("ignore_time_range", config.PlatformOptions.IgnoreTimeRange)
						h.WriteFloatAttribute("y_axis_left_min", config.YAxisLeft.Min)
						h.WriteFloatAttribute("y_axis_left_max", config.YAxisLeft.Max)
						h.WriteBlock("null_values", []string{}, func() {
							h.WriteStringAttribute("null_value", config.NullValues.NullValue)
							for _, so := range config.NullValues.SeriesOverrides {
								h.WriteBlock("series_overrides", []string{}, func() {
									h.WriteStringAttribute("null_value", so.NullValue)
									h.WriteStringAttribute("series_name", so.SeriesName)
								})
							}
						})
						h.WriteBlock("units", []string{}, func() {
							h.WriteStringAttribute("unit", config.Units.Unit)
							for _, so := range config.Units.SeriesOverrides {
								h.WriteBlock("series_overrides", []string{}, func() {
									h.WriteStringAttribute("unit", so.Unit)
									h.WriteStringAttribute("series_name", so.SeriesName)
								})
							}
						})
						h.WriteBlock("colors", []string{}, func() {
							h.WriteStringAttribute("color", config.Colors.Color)
							for _, so := range config.Colors.SeriesOverrides {
								h.WriteBlock("series_overrides", []string{}, func() {
									h.WriteStringAttribute("color", so.Color)
									h.WriteStringAttribute("series_name", so.SeriesName)
								})
							}
						})

					})
				}
			})
		}

		for _, v := range d.Variables {
			h.WriteBlock("variable", []string{}, func() {
				h.WriteStringAttribute("name", v.Name)
				h.WriteStringAttributeIfNotEmpty("title", v.Title)
				h.WriteStringAttribute("type", strings.ToLower(string(v.Type)))

				if v.DefaultValues != nil {
					arr := []string{}
					for _, value := range *v.DefaultValues {
						arr = append(arr, value.Value.String)
					}
					h.WriteStringSliceAttributeIfNotEmpty("default_values", arr)
				}

				if v.NRQLQuery != nil {
					h.WriteBlock("nrql_query", []string{}, func() {
						h.WriteIntArrayAttribute("account_ids", v.NRQLQuery.AccountIDs)
						h.WriteStringAttribute("query", string(v.NRQLQuery.Query))
					})
				}

				if v.Items != nil {
					for _, item := range v.Items {
						h.WriteBlock("item", []string{}, func() {
							h.WriteStringAttribute("value", item.Value)
							h.WriteStringAttributeIfNotEmpty("title", item.Title)
						})
					}
				}

				if v.IsMultiSelection {
					h.WriteBooleanAttribute("is_multi_selection", v.IsMultiSelection)
				}

				h.WriteStringAttributeIfNotEmpty("replacement_strategy", strings.ToLower(string(v.ReplacementStrategy)))

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
