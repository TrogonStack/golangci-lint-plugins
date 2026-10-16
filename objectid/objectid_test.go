package objectid_test

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/objectid"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	tests := map[string]struct {
		analyzer *analysis.Analyzer
		patterns []string
	}{
		"default": {
			analyzer: objectid.Analyzer,
			patterns: []string{"raw", "named", "generated", "implements", "example.com/store"},
		},
		"ignore names": {
			analyzer: newAnalyzer(t, objectid.Settings{IgnoreNames: []string{"projectID", "databaseID"}}),
			patterns: []string{"ignorenames"},
		},
		"suggestion": {
			analyzer: newAnalyzer(t, objectid.Settings{Suggestion: "internal/boid", UntypedConstants: true}),
			patterns: []string{"suggestion"},
		},
		"exclude packages": {
			analyzer: newAnalyzer(t, objectid.Settings{ExcludePackages: []objectid.PackagePrefix{"example.com/excluded"}}),
			patterns: []string{"example.com/excluded", "example.com/excluded/inner"},
		},
		"untyped constants": {
			analyzer: newAnalyzer(t, objectid.Settings{UntypedConstants: true}),
			patterns: []string{"untyped"},
		},
		"name patterns": {
			analyzer: newAnalyzer(t, objectid.Settings{NamePatterns: []string{"Ref"}}),
			patterns: []string{"custom"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			analysistest.Run(t, testdata, tt.analyzer, tt.patterns...)
		})
	}
}

func TestNewRejectsBadSettings(t *testing.T) {
	tests := map[string]struct {
		settings objectid.Settings
		want     error
	}{
		"empty name pattern": {
			settings: objectid.Settings{NamePatterns: []string{"ID", ""}},
			want:     objectid.ErrEmptyNamePattern,
		},
		"empty ignore name": {
			settings: objectid.Settings{IgnoreNames: []string{"projectID", ""}},
			want:     objectid.ErrEmptyIgnoreName,
		},
		"empty package prefix": {
			settings: objectid.Settings{ExcludePackages: []objectid.PackagePrefix{"example.com/pb", ""}},
			want:     objectid.ErrEmptyPackagePrefix,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := objectid.New(tt.settings)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func newAnalyzer(t *testing.T, settings objectid.Settings) *analysis.Analyzer {
	t.Helper()
	analyzer, err := objectid.New(settings)
	if err != nil {
		t.Fatal(err)
	}
	return analyzer
}
