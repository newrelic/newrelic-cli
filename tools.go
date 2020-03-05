// +build tools

package main

import (
	// build/test.mk
	_ "github.com/stretchr/testify/assert"

	// build/lint.mk
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/llorllale/go-gitlint/cmd/go-gitlint"
	_ "github.com/psampaz/go-mod-outdated"

	// build/document.mk
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "golang.org/x/tools/cmd/godoc"

	// build/release.mk
	_ "github.com/goreleaser/goreleaser"
)
