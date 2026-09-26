package staticerr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/staticerr"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	tests := map[string][]string{
		"package-level vars in every shape report nothing":       {"good/..."},
		"an errors.New inside a function is reported":            {"inline/..."},
		"a var declaration inside a function is reported":        {"localvar/..."},
		"a package-level var holding a func literal is reported": {"funclit/..."},
		"an errors.New handed to another call is reported":       {"nested/..."},
		"the errors package under another name is reported":      {"aliased/..."},
		"errors.New in a _test.go file is skipped":               {"testskip/..."},
	}

	for name, patterns := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, testdata, staticerr.Analyzer, patterns...)
		})
	}
}
