package semconvkey_test

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/semconvkey"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	configured, err := semconvkey.New(semconvkey.Settings{
		AllowedPackages: []semconvkey.PackagePrefix{"example.com/appsemconv"},
	})
	if err != nil {
		t.Fatal(err)
	}

	generated, err := semconvkey.New(semconvkey.Settings{
		AllowedPackages:      []semconvkey.PackagePrefix{"example.com/appsemconv"},
		GeneratedInstruments: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		analyzer *analysis.Analyzer
		patterns []string
	}{
		"default": {
			analyzer: semconvkey.Analyzer,
			patterns: []string{
				"literal",
				"localconst",
				"localcopy",
				"semconvok",
				"conversion",
				"compositelit",
				"metricname",
				"unconfigured",
				"generated",
				"testskip",
			},
		},
		"allowed packages": {
			analyzer: configured,
			patterns: []string{"configured", "example.com/appsemconv"},
		},
		"generated instruments": {
			analyzer: generated,
			patterns: []string{"handrolled", "example.com/appsemconv"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			analysistest.Run(t, testdata, tt.analyzer, tt.patterns...)
		})
	}
}

func TestNewRejectsEmptyPackagePrefix(t *testing.T) {
	_, err := semconvkey.New(semconvkey.Settings{
		AllowedPackages: []semconvkey.PackagePrefix{"example.com/appsemconv", ""},
	})
	if !errors.Is(err, semconvkey.ErrEmptyPackagePrefix) {
		t.Fatalf("got %v, want %v", err, semconvkey.ErrEmptyPackagePrefix)
	}
}

func TestNewRejectsPackagePrefixNotNamedForSemconv(t *testing.T) {
	_, err := semconvkey.New(semconvkey.Settings{
		AllowedPackages: []semconvkey.PackagePrefix{"example.com/appsemconv", "example.com/app/internal/telemetry"},
	})
	if !errors.Is(err, semconvkey.ErrPackagePrefixNotSemconv) {
		t.Fatalf("got %v, want %v", err, semconvkey.ErrPackagePrefixNotSemconv)
	}
}
