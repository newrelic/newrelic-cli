package ux

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type PlainProgress struct {
}

func NewPlainProgress() *PlainProgress {
	p := PlainProgress{}

	return &p
}

func (p *PlainProgress) Start(recipe types.Recipe) {
	msg := fmt.Sprintf("Installing %s", recipe.Name)

	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...\n")
}

func (p *PlainProgress) Success(recipe types.Recipe) {
	msg := fmt.Sprintf("Installing %s", recipe.Name)

	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...success.\n")
}

func (p *PlainProgress) Fail(recipe types.Recipe) {
	msg := fmt.Sprintf("Installing %s", recipe.Name)

	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...failed.\n")
}

func (p *PlainProgress) Stop() {}
