package connecterror_test

import (
	"errors"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TrogonStack/golangci-lint-plugins/connecterror"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	analyzer, err := connecterror.New(connecterror.Settings{Replacements: []string{"acme/orders/rpcerr.NewError"}})
	if err != nil {
		t.Fatal(err)
	}
	several, err := connecterror.New(connecterror.Settings{Replacements: []string{
		"acme/orders/rpcerr.NewError",
		"acme/orders/autherr.NewError",
		"acme/orders/rpcerr.NewError",
	}})
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

	t.Run("every replacement's package may call connect.NewError, and the diagnostic names them all", func(t *testing.T) {
		t.Parallel()
		analysistest.Run(t, analysistest.TestData(), several, "acme/orders/rpcerr/...", "acme/orders/autherr/...", "acme/orders/several/...")
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacements []string
		want         error
	}{
		"a function in a module path":         {[]string{"github.com/acme/orders/rpcerr.NewError"}, nil},
		"a function in a single-segment path": {[]string{"rpcerr.NewError"}, nil},
		"no replacement":                      {nil, connecterror.ErrReplacementRequired},
		"an empty entry":                      {[]string{""}, connecterror.ErrInvalidReplacement},
		"no function name":                    {[]string{"github.com/acme/orders/rpcerr"}, connecterror.ErrInvalidReplacement},
		"an unexported function":              {[]string{"github.com/acme/orders/rpcerr.newError"}, connecterror.ErrInvalidReplacement},
		"a method rather than a function":     {[]string{"github.com/acme/orders/rpcerr.Builder.NewError"}, connecterror.ErrInvalidReplacement},
		"no import path":                      {[]string{".NewError"}, connecterror.ErrInvalidReplacement},
		"several functions":                   {[]string{"github.com/acme/orders/rpcerr.NewError", "github.com/acme/orders/autherr.NewError"}, nil},
		"one bad function among several":      {[]string{"github.com/acme/orders/rpcerr.NewError", "github.com/acme/orders/autherr"}, connecterror.ErrInvalidReplacement},
		"connect.NewError itself":             {[]string{"connectrpc.com/connect.NewError"}, connecterror.ErrInvalidReplacement},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := connecterror.New(connecterror.Settings{Replacements: tt.replacements})
			if !errors.Is(err, tt.want) {
				t.Fatalf("New(%q) error = %v, want %v", tt.replacements, err, tt.want)
			}
		})
	}
}
