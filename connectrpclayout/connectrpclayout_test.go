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

	snakeCase, err := connectrpclayout.New(connectrpclayout.Settings{Naming: connectrpclayout.SnakeCase})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		analyzer *analysis.Analyzer
		patterns []string
	}{
		"a lowercase layout reports nothing by default":        {connectrpclayout.Analyzer, []string{"good/..."}},
		"a package handing a handler on is not an rpc package": {connectrpclayout.Analyzer, []string{"wiring/..."}},
		"an rpc package named for the wrong rpc is reported":   {connectrpclayout.Analyzer, []string{"badpkgname/..."}},
		"an rpc package at the wrong path is reported":         {connectrpclayout.Analyzer, []string{"wrongpath/..."}},
		"a handler under another name is reported":             {connectrpclayout.Analyzer, []string{"misnamed/..."}},
		"the handler declared in the wrong file is reported":   {connectrpclayout.Analyzer, []string{"wrongfile/..."}},
		"an rpc package declaring both forms is reported":      {connectrpclayout.Analyzer, []string{"bothforms/..."}},
		"declarations in _test.go files are skipped":           {connectrpclayout.Analyzer, []string{"testskip/..."}},
		"a snake_case layout reports nothing under SnakeCase":  {snakeCase, []string{"snake/..."}},
		"a lowercase layout is reported under SnakeCase":       {snakeCase, []string{"joined/..."}},
		"a package named apart from its directory is reported": {snakeCase, []string{"mixed/..."}},
		"a misplaced underscore is reported":                   {snakeCase, []string{"misspelled/..."}},
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
