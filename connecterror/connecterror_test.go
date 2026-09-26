package connecterror_test

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/connecterror"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	analyzer, err := connecterror.New(connecterror.Settings{Replacement: "acme/orders/rpcerr.NewError"})
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string][]string{
		"the replacement's own package may call connect.NewError": {"acme/orders/rpcerr/..."},
		"a direct call is reported":                               {"acme/orders/placeorder/..."},
		"a call through a renamed import is reported":             {"acme/orders/renamed/..."},
		"a call through a dot import is reported":                 {"acme/orders/dotimport/..."},
		"connect.NewError taken as a value is reported":           {"acme/orders/funcvalue/..."},
		"a call under Connect's previous import path is reported": {"acme/orders/legacy/..."},
		"a generated file is skipped":                             {"acme/orders/generated/..."},
		"a _test.go file is skipped":                              {"acme/orders/testsonly/..."},
		"a NewError of another package is not reported":           {"acme/orders/ownnewerror/..."},
	}

	for name, patterns := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			analysistest.Run(t, analysistest.TestData(), analyzer, patterns...)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacement string
		want        error
	}{
		"a function in a module path":         {"github.com/acme/orders/rpcerr.NewError", nil},
		"a function in a single-segment path": {"rpcerr.NewError", nil},
		"no replacement":                      {"", connecterror.ErrReplacementRequired},
		"no function name":                    {"github.com/acme/orders/rpcerr", connecterror.ErrInvalidReplacement},
		"an unexported function":              {"github.com/acme/orders/rpcerr.newError", connecterror.ErrInvalidReplacement},
		"a method rather than a function":     {"github.com/acme/orders/rpcerr.Builder.NewError", connecterror.ErrInvalidReplacement},
		"no import path":                      {".NewError", connecterror.ErrInvalidReplacement},
		"connect.NewError itself":             {"connectrpc.com/connect.NewError", connecterror.ErrInvalidReplacement},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := connecterror.New(connecterror.Settings{Replacement: tt.replacement})
			if !errors.Is(err, tt.want) {
				t.Fatalf("New(%q) error = %v, want %v", tt.replacement, err, tt.want)
			}
		})
	}
}
