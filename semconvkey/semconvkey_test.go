package semconvkey_test

import (
	"encoding/json"
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
				"dotimport",
				"alias",
				"methodexpr",
				"metricname",
				"unconfigured",
				"generated",
				"upstreamorigin",
				"attributereferences",
				"attributedotcapture",
				"factorycompat",
				"testskip",
			},
		},
		"allowed packages": {
			analyzer: configured,
			patterns: []string{"configured", "zeroliteral", "initialization", "example.com/appsemconv"},
		},
		"generated instruments": {
			analyzer: generated,
			patterns: []string{"handrolled", "references", "redeclared", "transitive", "attributereferences", "attributedotcapture", "initialization", "example.com/appsemconv"},
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

func TestGeneratedFilesCanBeChecked(t *testing.T) {
	var settings semconvkey.Settings
	if err := json.Unmarshal([]byte(`{"allowed-packages":["example.com/appsemconv"],"generated-instruments":true,"check-generated-files":true}`), &settings); err != nil {
		t.Fatal(err)
	}
	analyzer, err := semconvkey.New(settings)
	if err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), analyzer, "checkedgenerated", "example.com/appsemconvextra", "example.com/appsemconv")
}

func TestGeneratedFileCheckingKeepsNameOnlyMode(t *testing.T) {
	analyzer, err := semconvkey.New(semconvkey.Settings{CheckGeneratedFiles: true})
	if err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), analyzer, "checkedattribute", "factorycompat")
}
