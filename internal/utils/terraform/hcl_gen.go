package terraform

import (
	"fmt"
	"strings"
)

type HCLGen struct {
	strings.Builder
	shiftWidth int
	i          string
}

func NewHCLGen(shiftWidth int) *HCLGen {
	return &HCLGen{
		Builder:    strings.Builder{},
		shiftWidth: shiftWidth,
	}
}

func (h *HCLGen) WriteMultilineStringAttribute(label string, value string) {
	h.WriteString(fmt.Sprintf("%s%s = <<EOT\n%s\nEOT\n", h.i, label, value))
}

func (h *HCLGen) WriteStringAttribute(label string, value string) {
	h.WriteString(fmt.Sprintf("%s%s = \"%s\"\n", h.i, label, strings.ReplaceAll(value, "\"", "\\\"")))
}

func (h *HCLGen) WriteStringAttributeIfNotEmpty(label string, value string) {
	if value != "" {
		h.WriteStringAttribute(label, value)
	}
}

func (h *HCLGen) WriteStringSliceAttribute(label string, value []string) {
	h.WriteString(fmt.Sprintf("%s%s = [\"%s\"]", h.i, label, strings.Join(value, "\",\"")))
}

func (h *HCLGen) WriteStringSliceAttributeIfNotEmpty(label string, value []string) {
	if len(value) > 0 {
		h.WriteStringSliceAttribute(label, value)
	}
}

func (h *HCLGen) WriteIntAttribute(label string, value int) {
	h.WriteString(fmt.Sprintf("%s%s = %d\n", h.i, label, value))
}

func (h *HCLGen) WriteIntAttributeIfNotZero(label string, value int) {
	if value != 0 {
		h.WriteIntAttribute(label, value)
	}
}

func (h *HCLGen) WriteBlock(name string, labels []string, f func() error) error {
	h.WriteString(fmt.Sprintf("%s%s ", h.i, name))
	for _, l := range labels {
		h.WriteString(fmt.Sprintf("\"%s\" ", l))
	}
	h.WriteString("{\n")

	h.indent()
	if err := f(); err != nil {
		return err
	}
	h.unindent()

	h.WriteString(fmt.Sprintf("%s}\n", h.i))

	return nil
}

func (h *HCLGen) indent() {
	h.i += strings.Repeat(" ", h.shiftWidth)
}

func (h *HCLGen) unindent() {
	h.i = h.i[0 : len(h.i)-h.shiftWidth]
}
