package ctxflow_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/ctxflow"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ctxflow.Analyzer,
		"slogcalls",
		"siblings",
		"custom",
		"freshroot",
		"mainroot",
		"testfile",
		"generated",
	)
}
