// Package connecterror reports any use of connect.NewError outside the one
// function a codebase has chosen to build its Connect errors. Connect sends an
// error's Error() text to the client as is, so a codebase that decides what a
// client may read about a failure needs every Connect error to go through the
// function that makes that decision, and a direct connect.NewError is a way
// around it.
//
// The linter knows nothing about those functions beyond where they live: the
// Replacements setting names them, the diagnostic points at them, and the
// packages declaring them are the only ones allowed to call connect.NewError. Generated
// files and _test.go files are skipped.
package connecterror

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "reports connect.NewError outside the functions configured to build Connect errors"

var (
	ErrReplacementRequired = errors.New("at least one replacement is required")
	ErrInvalidReplacement  = errors.New("invalid replacement")
)

// connectPackages are the import paths Connect's Go runtime has shipped
// under: the current one, and the one it was published as before it moved.
var connectPackages = map[string]bool{
	"connectrpc.com/connect":         true,
	"github.com/bufbuild/connect-go": true,
}

const newErrorName = "NewError"

// Func is a package-level function, named by its import path and its name.
type Func struct {
	PkgPath string
	Name    string
}

// ParseFunc reads a function spelled as <import path>.<Name>, such as
// github.com/acme/orders/rpcerr.NewError.
func ParseFunc(s string) (Func, error) {
	slash := strings.LastIndex(s, "/")
	dot := strings.Index(s[slash+1:], ".")
	if dot < 0 {
		return Func{}, fmt.Errorf("%w %q: want <import path>.<Name>", ErrInvalidReplacement, s)
	}
	dot += slash + 1

	f := Func{PkgPath: s[:dot], Name: s[dot+1:]}
	if f.PkgPath == "" || !token.IsIdentifier(f.Name) || !token.IsExported(f.Name) {
		return Func{}, fmt.Errorf("%w %q: want <import path>.<Name> with an exported Name", ErrInvalidReplacement, s)
	}
	if connectPackages[f.PkgPath] {
		return Func{}, fmt.Errorf("%w %q: the replacement cannot be Connect itself", ErrInvalidReplacement, s)
	}
	return f, nil
}

func (f Func) String() string {
	return f.PkgPath + "." + f.Name
}

// Replacements are the functions a codebase builds its Connect errors with,
// in the order they were configured.
type Replacements []Func

// ParseReplacements reads each function with ParseFunc, keeping the first of
// any function named twice.
func ParseReplacements(ss []string) (Replacements, error) {
	if len(ss) == 0 {
		return nil, ErrReplacementRequired
	}
	var r Replacements
	seen := map[Func]bool{}
	for _, s := range ss {
		f, err := ParseFunc(s)
		if err != nil {
			return nil, err
		}
		if !seen[f] {
			seen[f] = true
			r = append(r, f)
		}
	}
	return r, nil
}

func (r Replacements) declaredIn(pkgPath string) bool {
	for _, f := range r {
		if f.PkgPath == pkgPath {
			return true
		}
	}
	return false
}

func (r Replacements) String() string {
	names := make([]string, len(r))
	for i, f := range r {
		names[i] = f.String()
	}
	if len(names) == 1 {
		return names[0]
	}
	return "one of " + strings.Join(names, ", ")
}

// Settings is what golangci-lint's custom linter settings decode into.
type Settings struct {
	Replacements []string `json:"replacements"`
}

// New is connecterror configured by settings. There is no default: which
// functions build a codebase's Connect errors is that codebase's decision.
func New(settings Settings) (*analysis.Analyzer, error) {
	replacements, err := ParseReplacements(settings.Replacements)
	if err != nil {
		return nil, err
	}
	return &analysis.Analyzer{
		Name: "connecterror",
		Doc:  Doc,
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, replacements)
		},
	}, nil
}

func run(pass *analysis.Pass, replacements Replacements) (any, error) {
	if replacements.declaredIn(pass.Pkg.Path()) {
		return nil, nil
	}

	for _, f := range pass.Files {
		if ast.IsGenerated(f) || strings.HasSuffix(pass.Fset.Position(f.Pos()).Filename, "_test.go") {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if ok && isConnectNewError(pass.TypesInfo.Uses[id]) {
				pass.Reportf(id.Pos(), "use %s instead of connect.NewError", replacements)
			}
			return true
		})
	}
	return nil, nil
}

func isConnectNewError(obj types.Object) bool {
	fn, ok := obj.(*types.Func)
	if !ok || fn.Name() != newErrorName || fn.Pkg() == nil {
		return false
	}
	return connectPackages[fn.Pkg().Path()] && fn.Type().(*types.Signature).Recv() == nil
}
