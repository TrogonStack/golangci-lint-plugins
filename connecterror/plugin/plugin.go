// Package plugin is how golangci-lint loads connecterror: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/connecterror"
)

func init() {
	register.Plugin("connecterror", newPlugin)
}

type linter struct {
	analyzer *analysis.Analyzer
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[connecterror.Settings](conf)
	if err != nil {
		return nil, err
	}

	analyzer, err := connecterror.New(settings)
	if err != nil {
		return nil, err
	}

	return linter{analyzer: analyzer}, nil
}

func (l linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{l.analyzer}, nil
}

// Types info, not syntax: the linter resolves an identifier to the function
// it names, so connect imported under another name, dot-imported, or passed
// as a value is still caught and somebody's own NewError is not.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
