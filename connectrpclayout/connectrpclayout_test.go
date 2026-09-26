package connectrpclayout_test

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/connectrpclayout"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	lowercase, err := connectrpclayout.New(connectrpclayout.Settings{Naming: connectrpclayout.Lowercase})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		analyzer *analysis.Analyzer
		patterns []string
	}{
		"a snake_case layout reports nothing by default":       {connectrpclayout.Analyzer, []string{"snake/..."}},
		"a lowercase layout is reported by default":            {connectrpclayout.Analyzer, []string{"joined/..."}},
		"a package named apart from its directory is reported": {connectrpclayout.Analyzer, []string{"mixed/..."}},
		"a misplaced underscore is reported":                   {connectrpclayout.Analyzer, []string{"misspelled/..."}},
		"a lowercase layout reports nothing under Lowercase":   {lowercase, []string{"good/..."}},
		"a package handing a handler on is not an rpc package": {lowercase, []string{"wiring/..."}},
		"an rpc package named for the wrong rpc is reported":   {lowercase, []string{"badpkgname/..."}},
		"an rpc package at the wrong path is reported":         {lowercase, []string{"wrongpath/..."}},
		"a handler under another name is reported":             {lowercase, []string{"misnamed/..."}},
		"the handler declared in the wrong file is reported":   {lowercase, []string{"wrongfile/..."}},
		"an rpc package declaring both forms is reported":      {lowercase, []string{"bothforms/..."}},
		"declarations in _test.go files are skipped":           {lowercase, []string{"testskip/..."}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, testdata, tt.analyzer, tt.patterns...)
		})
	}
}

func TestNewRejectsUnknownNaming(t *testing.T) {
	t.Parallel()

	_, err := connectrpclayout.New(connectrpclayout.Settings{Naming: "kebab-case"})
	if !errors.Is(err, connectrpclayout.ErrUnknownNaming) {
		t.Fatalf("got %v, want ErrUnknownNaming", err)
	}
}
