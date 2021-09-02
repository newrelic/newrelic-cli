package ux

import (
	"fmt"

	"github.com/fatih/color"
)

type PlainProgress struct {
}

func NewPlainProgress() *PlainProgress {
	p := PlainProgress{}

	return &p
}

func (p *PlainProgress) Start(msg string) {
	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...\n")
	fmt.Println()
}

func (p *PlainProgress) Success(msg string) {
	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...success.\n\n")
}

func (p *PlainProgress) Fail(msg string) {
	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...failed.\n\n")
}

func (p *PlainProgress) Canceled(msg string) {
	c := color.New(color.FgCyan)
	c.Printf("==>")
	x := color.New(color.Bold)
	x.Printf(" %s", msg)

	fmt.Printf("...canceled.\n\n")
}

func (p *PlainProgress) Stop() {}
