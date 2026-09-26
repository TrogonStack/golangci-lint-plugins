// Package plugin is how golangci-lint loads ctxflow: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/ctxflow"
)

func init() {
	register.Plugin("ctxflow", newPlugin)
}

type linter struct{}

// The rules have nothing to configure, so the settings are ignored.
func newPlugin(any) (register.LinterPlugin, error) {
	return linter{}, nil
}

func (linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{ctxflow.Analyzer}, nil
}

// Types info, not syntax: the linter resolves a callee and its Context
// sibling against their declared signatures, and resolves context.Background
// by what it is rather than how an import spells it.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
