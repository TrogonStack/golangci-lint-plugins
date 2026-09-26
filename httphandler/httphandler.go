// Package httphandler checks the convention for a plain http handler: it
// lives in its own package, and that package declares exactly one of Handler
// or NewHandler, in <packagename>.go, in the one shape each of those two
// spellings pins down.
//
// Handler is a func(http.ResponseWriter, *http.Request), for a handler that
// needs nothing from the process around it: a call site registers it with
// mux.HandleFunc. It is a func, not a var, because a package-level mutable
// handler is something any importer can reassign, and one spelling of the
// convention is the point. NewHandler is for a handler that does need
// something: it takes a single HandlerOptions declared beside it and returns
// an http.Handler, or an http.Handler and an error when it validates those
// options, so every handler package is wired the same way. A package
// declares one of the two, never both, so a package has one answer to what
// its handler is rather than two.
//
// Only a Handler or NewHandler whose declaration mentions a net/http type is
// in scope. A NewHandler that builds something other than an http handler,
// such as a log/slog handler, or a Connect rpc handler whose signature is
// spelled in generated types, is not what this convention governs.
package httphandler

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that a package declaring Handler or NewHandler declares exactly one of them, in the one shape each requires, in <packagename>.go"

// Analyzer is httphandler itself, for any analysis driver: golangci-lint
// through the plugin package, or a singlechecker binary.
var Analyzer = &analysis.Analyzer{
	Name:     "httphandler",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// The names this rule is checking against, spelled once here because the
// diagnostics quote them back to the reader.
const (
	handlerName    = "Handler"
	newHandlerName = "NewHandler"
	optionsName    = "HandlerOptions"
)

func run(pass *analysis.Pass) (any, error) {
	// package main is excluded: a command's entry point is not a package
	// anything else imports, so there is no importer for "one spelling of the
	// convention" to protect. This also excludes the "main" package the test
	// binary synthesizes for its own main (import path ending in ".test"),
	// which is not code anyone wrote.
	if pass.Pkg.Name() == "main" {
		return nil, nil
	}

	files := sortedFiles(pass)

	var nonTestFiles []*ast.File
	for _, f := range files {
		if !strings.HasSuffix(pass.Fset.Position(f.Pos()).Filename, "_test.go") {
			nonTestFiles = append(nonTestFiles, f)
		}
	}
	// No non-test files means this pass is the external xtest package for the
	// package (package foo_test compiled separately from package foo): its
	// handler, if it has one, lives in the sibling internal-package pass, not
	// this one.
	if len(nonTestFiles) == 0 {
		return nil, nil
	}

	handler := inScope(findTopLevelDecl(pass, nonTestFiles, handlerName))
	constructor := inScope(findTopLevelDecl(pass, nonTestFiles, newHandlerName))

	// Neither in scope means this package is not an http handler package: the
	// rule only governs a package that already claims to have one.
	if handler == nil && constructor == nil {
		return nil, nil
	}

	// Both is the ambiguous case rather than the generous one: it leaves two
	// answers to what this package's handler is, and an importer picks one,
	// so the package no longer says on its own which of them ships.
	if handler != nil && constructor != nil {
		pass.Reportf(constructor.ident.Pos(),
			"%s declares both %s and %s; a package has one answer to what its handler is, so declare %s for a handler that needs nothing from the process or %s for one that does",
			pass.Pkg.Name(), handlerName, newHandlerName, handlerName, newHandlerName)
		return nil, nil
	}

	if handler != nil {
		checkHandler(pass, handler)
		checkDeclFile(pass, handler)
		return nil, nil
	}

	checkConstructor(pass, constructor)
	checkDeclFile(pass, constructor)
	return nil, nil
}

// decl is one top-level, non-method declaration named Handler or NewHandler:
// either a func, or (the shape checkHandler rejects but this still has
// to recognize in order to reject it) a var.
type decl struct {
	ident *ast.Ident
	obj   types.Object
	fn    *ast.FuncDecl // nil when the declaration is a var, not a func
}

// checkHandler checks that Handler is a func(http.ResponseWriter,
// *http.Request), checked against net/http's actual types rather than
// against the identifiers ResponseWriter and Request, so a local type that
// happens to share one of those names does not satisfy it.
func checkHandler(pass *analysis.Pass, handler *decl) {
	fn, ok := handler.obj.(*types.Func)
	if !ok {
		pass.Reportf(handler.ident.Pos(), "%s must be func(http.ResponseWriter, *http.Request), got %s",
			handlerName, describeNonFunc(pass, handler.obj))
		return
	}

	sig := fn.Type().(*types.Signature)
	if isHandlerFuncSignature(sig) {
		return
	}
	pass.Reportf(handler.ident.Pos(), "%s must be func(http.ResponseWriter, *http.Request), got %s",
		handlerName, typeString(pass, sig))
}

// isHandlerFuncSignature reports whether sig is exactly
// func(http.ResponseWriter, *http.Request), with no result: the signature
// http.HandleFunc (and so mux.HandleFunc) accepts directly.
func isHandlerFuncSignature(sig *types.Signature) bool {
	params := sig.Params()
	return params.Len() == 2 && sig.Results().Len() == 0 &&
		isHTTPResponseWriter(params.At(0).Type()) &&
		isHTTPRequestPtr(params.At(1).Type())
}

// checkConstructor checks that NewHandler takes exactly one argument,
// a HandlerOptions declared beside it, and returns an http.Handler,
// optionally followed by an error for options it rejects. One struct parameter rather than a positional list, so a
// handler that gains an input later is a new field rather than a changed
// signature at every call site; a struct declared in this package rather than
// shared, because two handlers needing the same dependencies today is not a
// reason for one to break the other's signature tomorrow.
func checkConstructor(pass *analysis.Pass, constructor *decl) {
	fn, ok := constructor.obj.(*types.Func)
	if !ok {
		pass.Reportf(constructor.ident.Pos(), "%s must take exactly one %s declared in this package, got %s",
			newHandlerName, optionsName, describeNonFunc(pass, constructor.obj))
		return
	}
	sig := fn.Type().(*types.Signature)

	checkConstructorParams(pass, constructor.fn, sig)
	checkConstructorResults(pass, constructor.fn, sig)
}

func checkConstructorParams(pass *analysis.Pass, fn *ast.FuncDecl, sig *types.Signature) {
	params := sig.Params()
	if params.Len() != 1 {
		pass.Reportf(fn.Name.Pos(), "%s must take exactly one %s declared in this package, got %d parameters",
			newHandlerName, optionsName, params.Len())
		return
	}

	paramType := params.At(0).Type()
	named, ok := paramType.(*types.Named)
	if ok && named.Obj().Name() == optionsName && named.Obj().Pkg() == pass.Pkg {
		return
	}

	pass.Reportf(paramPos(fn), "the argument to %s must be a %s declared in this package, got %s",
		newHandlerName, optionsName, typeString(pass, paramType))
}

func checkConstructorResults(pass *analysis.Pass, fn *ast.FuncDecl, sig *types.Signature) {
	results := sig.Results()
	switch {
	case results.Len() == 1 && isHTTPHandler(results.At(0).Type()):
		return
	case results.Len() == 2 && isHTTPHandler(results.At(0).Type()) && isError(results.At(1).Type()):
		return
	}

	pass.Reportf(resultsPos(fn), "%s must return http.Handler or (http.Handler, error), got %s",
		newHandlerName, resultsString(pass, results))
}

// checkDeclFile checks that the declaration lives in
// <packagename>.go. A reader looking for a package's handler opens the file
// named after it.
func checkDeclFile(pass *analysis.Pass, d *decl) {
	wantFilename := pass.Pkg.Name() + ".go"
	gotFilename := filepath.Base(pass.Fset.Position(d.ident.Pos()).Filename)
	if gotFilename == wantFilename {
		return
	}
	pass.Reportf(d.ident.Pos(), "%s declares %s in %s, not %s", pass.Pkg.Name(), d.ident.Name, gotFilename, wantFilename)
}

// describeNonFunc renders what was declared instead of a func, for a
// diagnostic that says what was found. The var case is the one the rule
// explicitly calls out (a package-level mutable handler is reassignable by
// any importer), so it is named rather than folded into a generic fallback.
// The type case matters for a different reason: an identifier can be spelled
// Handler and still not be the handler this rule means, and saying so is how
// the check proves it looked at the declaration's kind rather than trusting
// its name.
func describeNonFunc(pass *analysis.Pass, obj types.Object) string {
	switch obj.(type) {
	case *types.Var:
		return "var of type " + typeString(pass, obj.Type())
	case *types.TypeName:
		return "a type declaration, not a func"
	default:
		return typeString(pass, obj.Type())
	}
}

// typeString renders t the way source in pass.Pkg would spell it: a type
// declared there stands alone, and a type from elsewhere is qualified by its
// package name, so a diagnostic says "otheropts.HandlerOptions" rather than
// the full import path.
func typeString(pass *analysis.Pass, t types.Type) string {
	return types.TypeString(t, func(pkg *types.Package) string {
		if pkg == pass.Pkg {
			return ""
		}
		return pkg.Name()
	})
}

// resultsString renders a func's results the way a reader would write them
// at a call site: a bare type for one result, a parenthesized list for any
// other count, so "got (error, http.Handler)" names the actual extra result
// rather than just the count.
func resultsString(pass *analysis.Pass, results *types.Tuple) string {
	n := results.Len()
	if n == 0 {
		return "no results"
	}
	if n == 1 {
		return typeString(pass, results.At(0).Type())
	}
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = typeString(pass, results.At(i).Type())
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

// paramPos and resultsPos point a diagnostic at the offending type rather
// than the func name, when there is one specific type to point at.
func paramPos(fn *ast.FuncDecl) token.Pos {
	if fn == nil || fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return token.NoPos
	}
	return fn.Type.Params.List[0].Type.Pos()
}

func resultsPos(fn *ast.FuncDecl) token.Pos {
	if fn == nil || fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
		return fn.Name.Pos()
	}
	return fn.Type.Results.List[0].Type.Pos()
}

// isNamedType reports whether t is the named type pkgPath.name, resolved
// through go/types rather than by matching the identifier's spelling alone:
// a local type that happens to share a name with a net/http type has a
// different Pkg, and does not match.
func isNamedType(t types.Type, pkgPath, name string) bool {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	pkg := obj.Pkg()
	return pkg != nil && pkg.Path() == pkgPath && obj.Name() == name
}

func isHTTPResponseWriter(t types.Type) bool {
	return isNamedType(t, "net/http", "ResponseWriter")
}

func isHTTPRequestPtr(t types.Type) bool {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	return ok && isNamedType(ptr.Elem(), "net/http", "Request")
}

func isError(t types.Type) bool {
	return types.Identical(t, types.Universe.Lookup("error").Type())
}

func isHTTPHandler(t types.Type) bool {
	return isNamedType(t, "net/http", "Handler")
}

// inScope is d when its declaration mentions net/http, and nil otherwise, so
// an unrelated declaration that shares the name is treated as absent.
func inScope(d *decl) *decl {
	if !mentionsHTTP(d) {
		return nil
	}
	return d
}

// mentionsHTTP reports whether d's declared type refers to a net/http type.
// A named type from elsewhere counts by its own package only, not by what it
// is built from, so a generated handler type whose request wraps http.Header
// deep inside does not pull its package into scope.
func mentionsHTTP(d *decl) bool {
	if d == nil {
		return false
	}
	t := d.obj.Type()
	if _, ok := d.obj.(*types.TypeName); ok {
		t = t.Underlying()
	}
	return refersToHTTP(t, map[types.Type]bool{})
}

func refersToHTTP(t types.Type, seen map[types.Type]bool) bool {
	if t == nil || seen[t] {
		return false
	}
	seen[t] = true

	switch t := t.(type) {
	case *types.Alias:
		return refersToHTTP(types.Unalias(t), seen)
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil && pkg.Path() == "net/http" {
			return true
		}
		for i := 0; i < t.TypeArgs().Len(); i++ {
			if refersToHTTP(t.TypeArgs().At(i), seen) {
				return true
			}
		}
		return false
	case *types.Pointer:
		return refersToHTTP(t.Elem(), seen)
	case *types.Slice:
		return refersToHTTP(t.Elem(), seen)
	case *types.Array:
		return refersToHTTP(t.Elem(), seen)
	case *types.Chan:
		return refersToHTTP(t.Elem(), seen)
	case *types.Map:
		return refersToHTTP(t.Key(), seen) || refersToHTTP(t.Elem(), seen)
	case *types.Signature:
		return refersToHTTP(t.Params(), seen) || refersToHTTP(t.Results(), seen)
	case *types.Tuple:
		for i := 0; i < t.Len(); i++ {
			if refersToHTTP(t.At(i).Type(), seen) {
				return true
			}
		}
		return false
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			if refersToHTTP(t.Field(i).Type(), seen) {
				return true
			}
		}
		return false
	case *types.Interface:
		for i := 0; i < t.NumEmbeddeds(); i++ {
			if refersToHTTP(t.EmbeddedType(i), seen) {
				return true
			}
		}
		for i := 0; i < t.NumExplicitMethods(); i++ {
			if refersToHTTP(t.ExplicitMethod(i).Type(), seen) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// findTopLevelDecl looks for a top-level, non-method declaration named name
// across files: a func, a var, or (so that an identifier merely spelled
// Handler cannot dodge the check by not being one) a type. There can be at
// most one: two top-level declarations sharing a name in the same package is
// a compile error, not something this analyzer needs to guard against.
func findTopLevelDecl(pass *analysis.Pass, files []*ast.File, name string) *decl {
	for _, f := range files {
		for _, top := range f.Decls {
			switch d := top.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && d.Name.Name == name {
					return &decl{ident: d.Name, obj: pass.TypesInfo.Defs[d.Name], fn: d}
				}
			case *ast.GenDecl:
				switch d.Tok {
				case token.VAR:
					for _, spec := range d.Specs {
						vs := spec.(*ast.ValueSpec)
						for _, ident := range vs.Names {
							if ident.Name == name {
								return &decl{ident: ident, obj: pass.TypesInfo.Defs[ident]}
							}
						}
					}
				case token.TYPE:
					for _, spec := range d.Specs {
						ts := spec.(*ast.TypeSpec)
						if ts.Name.Name == name {
							return &decl{ident: ts.Name, obj: pass.TypesInfo.Defs[ts.Name]}
						}
					}
				}
			}
		}
	}
	return nil
}

// sortedFiles returns pass.Files ordered by filename, so which file is
// picked first when scanning for a declaration does not depend on the
// driver's file ordering.
func sortedFiles(pass *analysis.Pass) []*ast.File {
	files := make([]*ast.File, len(pass.Files))
	copy(files, pass.Files)
	sort.Slice(files, func(i, j int) bool {
		return pass.Fset.Position(files[i].Pos()).Filename < pass.Fset.Position(files[j].Pos()).Filename
	})
	return files
}
