package connectrpclayout_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/connectrpclayout"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	tests := map[string][]string{
		"a concatenated layout reports nothing":                  {"good/..."},
		"a snake_case layout under another root reports nothing": {"snake/..."},
		"a package handing a handler on is not an rpc package":   {"wiring/..."},
		"an rpc package named for the wrong rpc is reported":     {"badpkgname/..."},
		"an rpc package at the wrong path is reported":           {"wrongpath/..."},
		"a handler under another name is reported":               {"misnamed/..."},
		"the handler declared in the wrong file is reported":     {"wrongfile/..."},
		"an rpc package declaring both forms is reported":        {"bothforms/..."},
		"declarations in _test.go files are skipped":             {"testskip/..."},
	}

	for name, patterns := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, testdata, connectrpclayout.Analyzer, patterns...)
		})
	}
}
