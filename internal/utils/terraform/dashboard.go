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
	ThresholdSeverityValues = map[string]string{
		"CRITICAL": "critical",
		"WARNING":  "warning",
	}
)

type DashboardBillboardSettings struct {
	Link        DashboardBillboardLink        `json:"link,omitempty"`
	Visual      DashboardBillboardVisual      `json:"visual,omitempty"`
	GridOptions DashboardBillboardGridOptions `json:"gridOptions,omitempty"`
}

type DashboardBillboardLink struct {
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
	NewTab bool   `json:"newTab,omitempty"`
}

type DashboardBillboardVisual struct {
	Alignment string `json:"alignment,omitempty"`
	Display   string `json:"display,omitempty"`
}

type DashboardBillboardGridOptions struct {
	Value   int `json:"value,omitempty"`
	Label   int `json:"label,omitempty"`
	Columns int `json:"columns,omitempty"`
}

type DashboardWidgetRawConfiguration struct {
	DataFormatters    []DataFormatter                `json:"dataFormatters"`
	NRQLQueries       []DashboardWidgetNRQLQuery     `json:"nrqlQueries"`
	LinkedEntityGUIDs []string                       `json:"linkedEntityGuids"`
	Text              string                         `json:"text"`
	Limit             float64                        `json:"limit,omitempty"`
	Facet             DashboardWidgetFacet           `json:"facet,omitempty"`
	Legend            DashboardWidgetLegend          `json:"legend,omitempty"`
	Threshold         json.RawMessage                `json:"thresholds,omitempty"`
	YAxisLeft         DashboardWidgetYAxisLeft       `json:"yAxisLeft,omitempty"`
	NullValues        DashboardWidgetNullValues      `json:"nullValues,omitempty"`
	Units             DashboardWidgetUnits           `json:"units,omitempty"`
	Colors            DashboardWidgetColors          `json:"colors,omitempty"`
	PlatformOptions   DashboardWidgetPlatformOptions `json:"platformOptions,omitempty"`
	RefreshRate       DashboardWidgetRefreshRate     `json:"refreshRate,omitempty"`
	InitialSorting    DashboardWidgetInitialSorting  `json:"initialSorting,omitempty"`
	BillboardSettings DashboardBillboardSettings   `json:"billboardSettings,omitempty"`
}

type DataFormatter struct {
	Name      string      `json:"name"`
	Precision interface{} `json:"precision"`
	Type      string      `json:"type"`
	Format    string      `json:"format"`
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

type DashboardWidgetThreshold struct {
	From     float64 `json:"from,omitempty"`
	Name     string  `json:"name,omitempty"`
	Severity string  `json:"severity,omitempty"`
	To       float64 `json:"to,omitempty"`
}

type DashboardWidgetLineThreshold struct {
	IsLabelVisible bool                       `json:"isLabelVisible,omitempty"`
	Threshold      []DashboardWidgetThreshold `json:"thresholds,omitempty"`
}

type DashboardWidgetBillBoardThreshold struct {
	AlertSeverity string  `json:"alertSeverity,omitempty"`
	Value         float64 `json:"value,omitempty"`
}

type DashboardWidgetYAxisLeft struct {
	Max  float64 `json:"max,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Zero bool    `json:"zero,omitempty"`
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

type DashboardWidgetRefreshRate struct {
	Frequency interface{} `json:"frequency,omitempty"`
}

type DashboardWidgetInitialSorting struct {
	Direction string `json:"direction"`
	Name      string `json:"name"`
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
						writeInterfaceValues(h, "refresh_rate", config.RefreshRate.Frequency) // function to handle different types of refresh rates which cant be handled through struct

						if w.Visualization.ID == "viz.line" {
							writeLineWidgetAttributes(h, config)
						}
						if w.Visualization.ID == "viz.billboard" {
							writeBillboardWidgetAttributes(h, config)
						}
						if w.Visualization.ID == "viz.table" {
							writeTableWidgetAttributes(h, config)
						}
						if w.Visualization.ID == "viz.bullet" {
							h.WriteFloatAttribute("limit", config.Limit)
						}

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

func writeLineWidgetAttributes(h *HCLGen, config *DashboardWidgetRawConfiguration) {
	h.WriteBooleanAttribute("y_axis_left_zero", config.YAxisLeft.Zero)
	var widgetLineThreshold DashboardWidgetLineThreshold
	if err := json.Unmarshal(config.Threshold, &widgetLineThreshold); err != nil {
		log.Fatal("Error unmarshalling widgetLineThreshold:", err)
	}

	h.WriteBooleanAttribute("is_label_visible", widgetLineThreshold.IsLabelVisible)

	for _, q := range widgetLineThreshold.Threshold {
		h.WriteBlock("threshold", []string{}, func() {
			h.WriteStringAttribute("name", q.Name)
			h.WriteStringAttribute("severity", q.Severity)
			h.WriteFloatAttribute("from", q.From)
			h.WriteFloatAttribute("to", q.To)
		})
	}
}

func writeBillboardWidgetAttributes(h *HCLGen, config *DashboardWidgetRawConfiguration) {
	var billboardThreshold []DashboardWidgetBillBoardThreshold
	if err := json.Unmarshal(config.Threshold, &billboardThreshold); err != nil {
		log.Fatal("Error unmarshalling billboardThreshold:", err)
	}
	for _, q := range billboardThreshold {
		h.WriteFloatAttribute(ThresholdSeverityValues[q.AlertSeverity], q.Value)
	}

	for _, q := range config.DataFormatters {
		h.WriteBlock("data_format", []string{}, func() {
			h.WriteStringAttribute("name", q.Name)
			h.WriteStringAttribute("type", q.Type)
			h.WriteStringAttribute("format", q.Format)
			writeInterfaceValues(h, "precision", q.Precision) // function to handle different types of precision
		})
	}

	if config.BillboardSettings != (DashboardBillboardSettings{}) {
		h.WriteBlock("billboard_settings", []string{}, func() {
			if config.BillboardSettings.Link != (DashboardBillboardLink{}) {
				h.WriteBlock("link", []string{}, func() {
					h.WriteStringAttributeIfNotEmpty("url", config.BillboardSettings.Link.URL)
					h.WriteStringAttributeIfNotEmpty("title", config.BillboardSettings.Link.Title)
					h.WriteBooleanAttribute("new_tab", config.BillboardSettings.Link.NewTab)
				})
			}

			if config.BillboardSettings.Visual != (DashboardBillboardVisual{}) {
				h.WriteBlock("visual", []string{}, func() {
					h.WriteStringAttributeIfNotEmpty("alignment", config.BillboardSettings.Visual.Alignment)
					h.WriteStringAttributeIfNotEmpty("display", config.BillboardSettings.Visual.Display)
				})
			}

			if config.BillboardSettings.GridOptions != (DashboardBillboardGridOptions{}) {
				h.WriteBlock("grid_options", []string{}, func() {
					h.WriteIntAttributeIfNotZero("columns", config.BillboardSettings.GridOptions.Columns)
					h.WriteIntAttributeIfNotZero("label", config.BillboardSettings.GridOptions.Label)
					h.WriteIntAttributeIfNotZero("value", config.BillboardSettings.GridOptions.Value)
				})
			}
		})
	}
}

func writeTableWidgetAttributes(h *HCLGen, config *DashboardWidgetRawConfiguration) {
	h.WriteBlock("initial_sorting", []string{}, func() {
		h.WriteStringAttribute("name", config.InitialSorting.Name)
		h.WriteStringAttribute("direction", config.InitialSorting.Direction)
	})

	for _, q := range config.DataFormatters {
		h.WriteBlock("data_format", []string{}, func() {
			h.WriteStringAttribute("name", q.Name)
			h.WriteStringAttribute("type", q.Type)
			h.WriteStringAttribute("format", q.Format)
			writeInterfaceValues(h, "precision", q.Precision) // function to handle different types of precision
		})
	}
}

func writeInterfaceValues(h *HCLGen, title string, titleValue interface{}) {
	switch titleValue := titleValue.(type) {
	case string:
		h.WriteStringAttribute(title, titleValue) // string without quotes
	case float64:
		h.WriteFloatAttribute(title, titleValue) // integer without quotes
	default:
		h.WriteStringAttribute(title, "")
	}
}
