// Package staticerr checks that an error with a fixed message is declared
// once, at package level, rather than built where it is returned.
//
// It reports a call to errors.New anywhere but the right-hand side of a
// package-level var, skipping *_test.go files, where a one-off error is a
// fixture standing in for a dependency's failure and has no caller to match
// it.
//
// The reason is that errors.New of a constant string is a value whose identity
// is stable and costs nothing to keep. Calling it where the error is returned
// throws that identity away: one condition becomes as many distinct,
// non-equal error values as there are calls, and errors.Is can match none of
// them. A caller is left comparing message text, which is the thing wrapping
// exists to stop, or it gives up and treats every failure of the function
// alike. Hoisting the call to a package-level var costs one line and one name,
// and the name is the part that matters: it is what makes the condition
// something a caller can ask about.
//
// The rule is deliberately about errors.New and not about fmt.Errorf, which is
// the other common way to build an error. An fmt.Errorf interpolates, so there
// is no one value to hoist and nothing to give a stable identity to; its
// answer is to wrap a sentinel with %w, which is errorlint's job and not this
// linter's. The asymmetry is not an oversight. It follows from which of the
// two has an identity to lose.
package staticerr

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that errors.New is called only in a package-level var, so a fixed error has one identity"

// Analyzer is staticerr itself, for any analysis driver: golangci-lint
// through staticerr/plugin, or a singlechecker binary.
var Analyzer = &analysis.Analyzer{
	Name:     "staticerr",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		checkCall(pass, n.(*ast.CallExpr), stack)

		return true
	})

	return nil, nil
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr, stack []ast.Node) {
	if strings.HasSuffix(pass.Fset.Position(call.Pos()).Filename, "_test.go") {
		return
	}

	if !isErrorsNew(pass, call) || inPackageLevelVar(stack) {
		return
	}

	pass.Reportf(call.Pos(),
		"errors.New(%s) is built where it is returned, so every call is a different error and errors.Is can match none of them; declare it once as a package-level var and return that",
		messageOf(call))
}

// isErrorsNew resolves the callee through pass.TypesInfo rather than matching
// the syntax "errors.New", so a package aliased under another name is caught
// and a method of somebody's own type called New is not.
func isErrorsNew(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	fn, ok := pass.TypesInfo.Uses[sel.Sel].(*types.Func)
	if !ok || fn.Pkg() == nil {
		return false
	}

	return fn.Pkg().Path() == "errors" && fn.Name() == "New"
}

// inPackageLevelVar reports whether the call sits in a var declaration at file
// scope. A function is what disqualifies it, and the check for one comes
// first: a package-level var may hold a func literal, and the var declaration
// is above that literal on the stack, so a search that stopped at the first
// var would exempt a call the literal makes at run time rather than once at
// initialisation.
func inPackageLevelVar(stack []ast.Node) bool {
	found := false

	for _, node := range stack {
		switch decl := node.(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			return false
		case *ast.GenDecl:
			found = found || decl.Tok == token.VAR
		}
	}

	return found
}

// messageOf is the call's first argument as written, so the report names the
// error it is about. A call whose message is not a literal is described by
// what it is rather than quoted, since the point of the report is the same
// either way.
func messageOf(call *ast.CallExpr) string {
	if len(call.Args) != 1 {
		return "..."
	}

	if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return lit.Value
	}

	return "..."
}
