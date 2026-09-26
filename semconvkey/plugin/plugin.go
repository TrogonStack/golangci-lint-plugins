// Package plugin is how golangci-lint loads semconvkey: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/semconvkey"
)

func init() {
	register.Plugin("semconvkey", newPlugin)
}

type linter struct {
	analyzer *analysis.Analyzer
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[semconvkey.Settings](conf)
	if err != nil {
		return nil, err
	}

	analyzer, err := semconvkey.New(settings)
	if err != nil {
		return nil, err
	}

	return linter{analyzer: analyzer}, nil
}

func (l linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{l.analyzer}, nil
}

// Types info, not syntax: the linter resolves a call's callee and a name's
// declaration to the package they come from, so that the attribute package
// under another name is still caught and somebody's own method of the same
// name is not.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
