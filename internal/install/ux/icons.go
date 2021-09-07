package ux

import (
	"github.com/fatih/color"
)

var (
	IconCheckmark      = "\u2705"
	IconMultiplication = "\u274C"
	IconMinus          = "\u2796"
	IconArrowRight     = "\u2B95"

	IconSuccess     = color.GreenString(IconCheckmark)
	IconError       = color.RedString(IconMultiplication)
	IconUnsupported = "🚫"
)
