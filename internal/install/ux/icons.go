package ux

import (
	"github.com/fatih/color"
)

// Unicode characters
// https://unicode-table.com/en/
var (
	IconCheckmark      = "\u2714"
	IconMultiplication = "\u274C"
	IconMinus          = "\u2212"
	IconArrowRight     = "\u2B95"
	IconExclamation    = "\u0021"
	IconCircleSlash    = "\u2298"

	IconSuccess     = color.GreenString(IconCheckmark)
	IconError       = color.YellowString(IconExclamation) // We display "warning"	 symbol to avoid scary "red" colors
	IconUnsupported = color.RedString(IconCircleSlash)
)
