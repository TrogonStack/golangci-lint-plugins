// Package plugin is how golangci-lint loads connectrpclayout: it can only run
// a linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/connectrpclayout"
)

func init() {
	register.Plugin("connectrpclayout", newPlugin)
}

type linter struct{}

// The layout has no knobs, so there is nothing to configure and the settings
// are ignored.
func newPlugin(any) (register.LinterPlugin, error) {
	return linter{}, nil
}

func (linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{connectrpclayout.Analyzer}, nil
}

// Types info, not syntax: the linter recognises an rpc package by the
// generated handler types it declares, which syntax alone does not know.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
