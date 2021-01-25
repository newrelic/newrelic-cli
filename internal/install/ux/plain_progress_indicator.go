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
	fmt.Printf(" %s... ", msg)
}

func (p *PlainProgress) Success() {
	fmt.Println("success")
}

func (p *PlainProgress) Fail() {
	fmt.Println("fail")
}

func (p *PlainProgress) Stop() {}
