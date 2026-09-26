package httphandler_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/httphandler"
)

// TestAnalyzer checks every case in one analysistest.Run: the testdata
// imports net/http, and a Run per case would type-check net/http per case.
// The case names are kept as keys, so a reader still sees what each package
// is for, and a failure still names the testdata file and line it came from.
func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	// Each key says what its package proves; the analyzer sees them together.
	tests := map[string][]string{
		"a compliant Handler package reports nothing":                            {"good/..."},
		"a compliant NewHandler package reports nothing":                         {"goodnewhandler/..."},
		"a package declaring both forms is reported":                             {"bothforms/..."},
		"Handler with the wrong signature is reported":                           {"badsig/..."},
		"Handler declared as a var is reported":                                  {"handlervar/..."},
		"NewHandler with two parameters is reported":                             {"twoparams/..."},
		"NewHandler taking a struct not named HandlerOptions is reported":        {"badoptsname/..."},
		"NewHandler taking a HandlerOptions imported from elsewhere is reported": {"importedopts/..."},
		"NewHandler returning (http.Handler, error) is reported":                 {"badreturn/..."},
		"NewHandler returning a type other than http.Handler is reported":        {"badreturntype/..."},
		"the declaration in the wrong file is reported":                          {"wrongfile/..."},
		"a NewHandler whose signature never mentions net/http is skipped":        {"sloglike/...", "connectlike/..."},
		"a Handler in package main is skipped":                                   {"mainskip/..."},
		"a Handler declared only in a _test.go file is skipped":                  {"testonly/..."},
		"a local interface named Handler is reported as a bad Handler":           {"localiface/..."},
	}

	patterns := make([]string, 0, len(tests))
	for _, p := range tests {
		patterns = append(patterns, p...)
	}
	analysistest.Run(t, testdata, httphandler.Analyzer, patterns...)
}
