// Package plugin is how golangci-lint loads objectid: it can only run a
// linter that registered itself in a binary built from .custom-gcl.yml, and
// this is that registration. It is kept apart from the analyzer so that using
// the analyzer on its own does not build against golangci-lint's plugin API.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TrogonStack/golangci-lint-plugins/objectid"
)

func init() {
	register.Plugin("objectid", newPlugin)
}

type linter struct {
	analyzer *analysis.Analyzer
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[objectid.Settings](conf)
	if err != nil {
		return nil, err
	}

	analyzer, err := objectid.New(settings)
	if err != nil {
		return nil, err
	}

	return linter{analyzer: analyzer}, nil
}

func (l linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{l.analyzer}, nil
}

// Types info, not syntax: the linter tells type OrderID string from string
// and an alias of string by what the type is rather than how it is spelled,
// and matches a method against the interfaces of the packages it imports.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
