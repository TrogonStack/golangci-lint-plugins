// Package plugin is how golangci-lint loads httphandler: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/httphandler"
)

func init() {
	register.Plugin("httphandler", newPlugin)
}

type linter struct{}

// The convention has nothing to configure, so the settings are ignored.
func newPlugin(any) (register.LinterPlugin, error) {
	return linter{}, nil
}

func (linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{httphandler.Analyzer}, nil
}

// Types info, not syntax: the linter resolves a declaration's type against
// net/http's actual types, and a type is not something syntax alone knows.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
