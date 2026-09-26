// Package plugin is how golangci-lint loads staticerr: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/staticerr"
)

func init() {
	register.Plugin("staticerr", newPlugin)
}

type linter struct{}

// The rule has no knobs, so there is nothing to configure and the settings
// are ignored.
func newPlugin(any) (register.LinterPlugin, error) {
	return linter{}, nil
}

func (linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{staticerr.Analyzer}, nil
}

// Types info, not syntax: the linter resolves a call's callee to the package
// and function it names, so that errors imported under another name is still
// caught and somebody's own New is not.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
