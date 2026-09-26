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

type linter struct {
	analyzer *analysis.Analyzer
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[connectrpclayout.Settings](conf)
	if err != nil {
		return nil, err
	}

	analyzer, err := connectrpclayout.New(settings)
	if err != nil {
		return nil, err
	}

	return linter{analyzer: analyzer}, nil
}

func (l linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{l.analyzer}, nil
}

// Types info, not syntax: the linter recognises an rpc package by the
// generated handler types it declares, which syntax alone does not know.
func (linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
