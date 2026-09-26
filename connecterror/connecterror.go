// Package connecterror reports any use of connect.NewError outside the one
// function a codebase has chosen to build its Connect errors. Connect sends an
// error's Error() text to the client as is, so a codebase that decides what a
// client may read about a failure needs every Connect error to go through the
// function that makes that decision, and a direct connect.NewError is a way
// around it.
//
// The linter knows nothing about that function beyond where it lives: the
// Replacement setting names it, the diagnostic points at it, and the package
// declaring it is the one package allowed to call connect.NewError. Generated
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
const Doc = "reports connect.NewError outside the function configured to build Connect errors"

var (
	ErrReplacementRequired = errors.New("replacement is required")
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

// Settings is what golangci-lint's custom linter settings decode into.
type Settings struct {
	Replacement string `json:"replacement"`
}

// New is connecterror configured by settings. There is no default: which
// function builds a codebase's Connect errors is that codebase's decision.
func New(settings Settings) (*analysis.Analyzer, error) {
	if settings.Replacement == "" {
		return nil, ErrReplacementRequired
	}
	replacement, err := ParseFunc(settings.Replacement)
	if err != nil {
		return nil, err
	}
	return &analysis.Analyzer{
		Name: "connecterror",
		Doc:  Doc,
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, replacement)
		},
	}, nil
}

func run(pass *analysis.Pass, replacement Func) (any, error) {
	if pass.Pkg.Path() == replacement.PkgPath {
		return nil, nil
	}

	for _, f := range pass.Files {
		if ast.IsGenerated(f) || strings.HasSuffix(pass.Fset.Position(f.Pos()).Filename, "_test.go") {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if ok && isConnectNewError(pass.TypesInfo.Uses[id]) {
				pass.Reportf(id.Pos(), "use %s instead of connect.NewError", replacement)
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
