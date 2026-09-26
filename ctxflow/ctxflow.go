// Package ctxflow checks that a context.Context, once a program has one,
// keeps flowing down every call instead of being dropped along the way.
//
// It reports:
//
//   - a call to a function or method F when F takes no context.Context but
//     the same package, or the same receiver, declares FContext or
//     FWithContext taking a context.Context first and otherwise the exact
//     same parameters and results, such as slog.Info next to
//     slog.InfoContext, (*sql.DB).Query next to QueryContext, or
//     http.NewRequest next to NewRequestWithContext;
//   - any use of context.Background or context.TODO outside the roots of a
//     program: func main in package main, a func init, TestMain, and an
//     Example function, not counting a function literal inside one.
//
// The reason is that a context carries what a call is part of: its
// deadline, its cancellation, and the trace and request values an
// observability handler reads back out of it. A log record written without
// it is not attached to the span it happened under, and a query or a
// request sent without it outlives the caller that gave up on it. The
// sibling rule is inferred from the shape of an API rather than from a list
// of loggers, so any library following the FContext convention is covered
// without configuration.
//
// A fresh root is only ever the right answer where no context exists yet,
// which is what the roots of a program are. Everywhere else there is one to
// pass: a test has t.Context(), and work that has to outlive its caller's
// cancellation, such as a flush during shutdown, keeps its values with
// context.WithoutCancel(ctx). A callback a library invokes with no context
// at all is the one place neither exists, and is what a //nolint with a
// reason is for. Generated files are skipped, since they are never
// hand-edited to begin with.
package ctxflow

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that a context.Context is passed down instead of dropped by a call that has a Context variant or replaced by context.Background or context.TODO"

// Analyzer is ctxflow itself, for any analysis driver: golangci-lint
// through the plugin package, or a singlechecker binary.
var Analyzer = &analysis.Analyzer{
	Name: "ctxflow",
	Doc:  Doc,
	Run:  run,
}

const contextPackage = "context"

// contextSuffixes are the spellings a context-taking sibling of F goes by
// in the standard library and the libraries that follow it: InfoContext,
// QueryContext and CommandContext, or NewRequestWithContext.
var contextSuffixes = []string{"Context", "WithContext"}

// freshRoots are the context package functions that start a context from
// nothing rather than derive one from the caller's.
var freshRoots = map[string]bool{
	"Background": true,
	"TODO":       true,
}

const (
	siblingMessage   = "%s drops the context; call %s with the caller's ctx instead"
	freshRootMessage = "context.%s starts a new context and loses the caller's deadline, cancellation and trace; pass ctx down, use t.Context() in a test, or context.WithoutCancel(ctx) for work that must outlive its caller"
)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		testFile := strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, "_test.go")
		for _, decl := range file.Decls {
			inspect(pass, decl, isRoot(pass, decl, testFile))
		}
	}
	return nil, nil
}

// inspect checks node, exempting it from the fresh-root rule when root is
// set. A function literal inside a root is checked as if it were not one,
// since it can capture the context its root has already built.
func inspect(pass *analysis.Pass, node ast.Node, root bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			if root {
				inspect(pass, n, false)
				return false
			}
		case *ast.CallExpr:
			checkSibling(pass, n)
		case *ast.Ident:
			if !root {
				checkFreshRoot(pass, n)
			}
		}
		return true
	})
}

// isRoot reports whether decl is where a program starts with no context to
// inherit: func main in package main, a func init in any package, or a
// test binary's TestMain or Example function, which unlike a test or a
// benchmark are given no *testing.T or *testing.B to ask for one.
func isRoot(pass *analysis.Pass, decl ast.Decl, testFile bool) bool {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok || fn.Recv != nil {
		return false
	}
	switch fn.Name.Name {
	case "init":
		return true
	case "main":
		return pass.Pkg.Name() == "main"
	case "TestMain":
		return testFile
	}
	return testFile && strings.HasPrefix(fn.Name.Name, "Example")
}

// checkFreshRoot reports ident when it names context.Background or
// context.TODO. It checks the name rather than the call so that a function
// value, such as newCtx := context.Background, cannot carry one past the
// rule.
func checkFreshRoot(pass *analysis.Pass, ident *ast.Ident) {
	fn, ok := pass.TypesInfo.Uses[ident].(*types.Func)
	if !ok || fn.Pkg() == nil || fn.Pkg().Path() != contextPackage || !freshRoots[fn.Name()] {
		return
	}
	pass.Reportf(ident.Pos(), freshRootMessage, fn.Name())
}

func checkSibling(pass *analysis.Pass, call *ast.CallExpr) {
	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return
	}
	sig := fn.Type().(*types.Signature)
	if takesContext(sig) {
		return
	}
	recv := receiverAt(pass, call)
	if sig.Recv() != nil && recv.typ != nil {
		// Resolved through the receiver as written, a method of a generic type
		// comes back instantiated, the same as the sibling lookup returns it.
		obj, _, _ := types.LookupFieldOrMethod(recv.typ, true, fn.Pkg(), fn.Name())
		if self, ok := obj.(*types.Func); ok {
			sig = self.Type().(*types.Signature)
		}
	}
	sibling := findSibling(fn, sig, recv)
	if sibling == nil {
		return
	}
	pass.Reportf(call.Lparen, siblingMessage, displayName(pass, fn), sibling.Name())
}

// findSibling returns the FContext or FWithContext that fn's package or
// receiver declares beside it, when it is fn with a context.Context
// prepended and nothing else changed. A sibling whose remaining parameters
// or results differ is a different operation that happens to share a
// prefix, such as signal.Notify and signal.NotifyContext, and is not one fn
// could be swapped for.
func findSibling(fn *types.Func, sig *types.Signature, recv receiver) *types.Func {
	for _, suffix := range contextSuffixes {
		candidate := lookup(fn, sig, recv, fn.Name()+suffix)
		if candidate == nil {
			continue
		}
		candidateSig := instantiateLike(sig, candidate.Type().(*types.Signature))
		if candidateSig != nil && isContextVariant(sig, candidateSig) {
			return candidate
		}
	}
	return nil
}

// receiver is the operand a method is called on, as the call site has it.
type receiver struct {
	typ         types.Type
	addressable bool
}

// receiverAt returns the receiver of call as written, so that a sibling is
// only suggested when it could be called on that same operand: a
// pointer-receiver LoadContext is no replacement for Load on a value that is
// not addressable, such as the result of a call.
func receiverAt(pass *analysis.Pass, call *ast.CallExpr) receiver {
	fun := ast.Unparen(call.Fun)
	switch index := fun.(type) {
	case *ast.IndexExpr:
		fun = ast.Unparen(index.X)
	case *ast.IndexListExpr:
		fun = ast.Unparen(index.X)
	}
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return receiver{}
	}
	selection, ok := pass.TypesInfo.Selections[sel]
	if !ok {
		return receiver{}
	}
	return receiver{
		typ:         selection.Recv(),
		addressable: selection.Kind() == types.MethodVal && pass.TypesInfo.Types[sel.X].Addressable(),
	}
}

// lookup returns the function or method named name beside fn. A method is
// only accepted when it is declared on fn's own receiver type, so that a
// method promoted from an embedded field is never offered as the sibling of
// one the outer type declares itself.
func lookup(fn *types.Func, sig *types.Signature, recv receiver, name string) *types.Func {
	if sig.Recv() == nil {
		candidate, _ := fn.Pkg().Scope().Lookup(name).(*types.Func)
		return candidate
	}
	if recv.typ == nil {
		return nil
	}
	obj, _, _ := types.LookupFieldOrMethod(recv.typ, recv.addressable, fn.Pkg(), name)
	candidate, ok := obj.(*types.Func)
	if !ok || !types.Identical(receiverBase(fn), receiverBase(candidate)) {
		return nil
	}
	return candidate
}

func receiverBase(fn *types.Func) types.Type {
	t := fn.Type().(*types.Signature).Recv().Type()
	if pointer, ok := t.(*types.Pointer); ok {
		t = pointer.Elem()
	}
	if named, ok := types.Unalias(t).(*types.Named); ok {
		return named.Origin()
	}
	return t
}

// instantiateLike returns candidate with its type parameters replaced by
// sig's, so that a generic F and FContext compare equal when they are the
// same shape. Each declares its own type parameters, which are never
// identical to one another as they stand. It returns nil when the two do not
// declare the same number of them.
func instantiateLike(sig, candidate *types.Signature) *types.Signature {
	tparams := sig.TypeParams()
	if tparams.Len() != candidate.TypeParams().Len() {
		return nil
	}
	if tparams.Len() == 0 {
		return candidate
	}
	args := make([]types.Type, 0, tparams.Len())
	for tparam := range tparams.TypeParams() {
		args = append(args, tparam)
	}
	instance, err := types.Instantiate(nil, candidate, args, false)
	if err != nil {
		return nil
	}
	return instance.(*types.Signature)
}

func isContextVariant(sig, candidate *types.Signature) bool {
	params, candidateParams := sig.Params(), candidate.Params()
	if candidateParams.Len() != params.Len()+1 || !isContext(candidateParams.At(0).Type()) {
		return false
	}
	if sig.Variadic() != candidate.Variadic() {
		return false
	}
	for i := range params.Len() {
		if !types.Identical(params.At(i).Type(), candidateParams.At(i+1).Type()) {
			return false
		}
	}
	return types.Identical(sig.Results(), candidate.Results())
}

func takesContext(sig *types.Signature) bool {
	for param := range sig.Params().Variables() {
		if isContext(param.Type()) {
			return true
		}
	}
	return false
}

func isContext(t types.Type) bool {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == contextPackage && obj.Name() == "Context"
}

// displayName spells fn the way the reader's own code would, package
// qualified for a function and receiver qualified for a method, so the
// diagnostic reads as slog.Info or (*slog.Logger).Info.
func displayName(pass *analysis.Pass, fn *types.Func) string {
	qualifier := func(p *types.Package) string {
		if p == pass.Pkg {
			return ""
		}
		return p.Name()
	}
	sig := fn.Type().(*types.Signature)
	if recv := sig.Recv(); recv != nil {
		return "(" + types.TypeString(recv.Type(), qualifier) + ")." + fn.Name()
	}
	if prefix := qualifier(fn.Pkg()); prefix != "" {
		return prefix + "." + fn.Name()
	}
	return fn.Name()
}
