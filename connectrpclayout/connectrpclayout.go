// Package connectrpclayout checks where Connect rpc handlers live, for the
// service structs protoc-gen-connect-go-servicestruct generates: each handler
// sits in a package of its own at <service>/<rpc>, named for the rpc, and
// declares it in the file named for the package as Handler or NewHandler. A
// reader looking for how an rpc is answered then has one place to open, and
// the wiring that sets a service struct's fields only ever has two names to
// look for.
//
// An rpc package is recognised by its types rather than by where it sits: it
// is a package that declares a handler whose type is one of the generated
// struct's handler func types. That keeps the rule free of any fixed root
// directory, and leaves helper packages beside the rpc packages alone. The
// path and package name spell the service and rpc in one Naming, snake_case
// unless configured otherwise: echo_service/echo_stream for the EchoStream rpc
// of EchoService, or echoservice/echostream under Lowercase.
//
// An rpc package declares exactly one of the two. Handler is the handler
// itself, for an rpc that needs nothing from the process around it.
// NewHandler is for one that does, and returns the handler; what it takes is
// the package's own business.
package connectrpclayout

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that Connect rpc handlers live in <service>/<rpc> packages named for the rpc and declare Handler or NewHandler in the file named for the package"

// Naming is how a service or rpc name is spelled as a directory and a package
// name.
type Naming string

const (
	// SnakeCase splits words with underscores: echo_service/echo_stream. It is
	// the default, because a name of several words run together is hard to
	// read back.
	SnakeCase Naming = "snake_case"
	// Lowercase runs the words together: echoservice/echostream, the spelling
	// Go's own package naming advice asks for.
	Lowercase Naming = "lowercase"
)

var ErrUnknownNaming = errors.New("unknown naming")

// orDefault is n, or SnakeCase when n is unset, or an error when n is neither
// of the namings the linter knows.
func (n Naming) orDefault() (Naming, error) {
	switch n {
	case "":
		return SnakeCase, nil
	case SnakeCase, Lowercase:
		return n, nil
	default:
		return "", fmt.Errorf("%w %q: want %q or %q", ErrUnknownNaming, n, SnakeCase, Lowercase)
	}
}

func (n Naming) spell(name string) string {
	if n == Lowercase {
		return strings.ToLower(name)
	}
	return snakeCase(name)
}

// Settings is what golangci-lint's custom linter settings decode into.
type Settings struct {
	Naming Naming `json:"naming"`
}

// Analyzer is connectrpclayout with the default settings, for any analysis
// driver: golangci-lint through connectrpclayout/plugin, or a singlechecker
// binary.
var Analyzer = newAnalyzer(SnakeCase)

// New is connectrpclayout configured by settings.
func New(settings Settings) (*analysis.Analyzer, error) {
	naming, err := settings.Naming.orDefault()
	if err != nil {
		return nil, err
	}
	return newAnalyzer(naming), nil
}

func newAnalyzer(naming Naming) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "connectrpclayout",
		Doc:  Doc,
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, naming)
		},
	}
}

func run(pass *analysis.Pass, naming Naming) (any, error) {
	// A driver that loads test variants (analysistest, and go vet's own driver
	// with -test) synthesizes a package for the test binary's main, named
	// "main" with an import path ending in ".test". It is not code anyone
	// wrote, so it is not what any of these rules is checking.
	if pass.Pkg.Name() == "main" && strings.HasSuffix(pass.Pkg.Path(), ".test") {
		return nil, nil
	}

	c := &checker{pass: pass, naming: naming, services: map[*types.Package][]rpcField{}}
	c.checkHandlerPackage()
	return nil, nil
}

// The names an rpc package is checked against. They are spelled once here
// because the diagnostics quote them back to the reader, and a message naming
// a different identifier than the one the rule enforces is worse than no
// message.
const (
	handlerName    = "Handler"
	newHandlerName = "NewHandler"
)

// rpcField is one <Rpc>Func field of a generated service struct.
type rpcField struct {
	service string     // e.g. "EchoService"
	rpc     string     // e.g. "EchoStream"
	typ     types.Type // the field's type, e.g. EchoServiceEchoStreamHandlerFunc
}

type checker struct {
	pass     *analysis.Pass
	naming   Naming
	services map[*types.Package][]rpcField
}

// path is where the package serving field's rpc lives.
func (c *checker) path(field rpcField) string {
	return c.naming.spell(field.service) + "/" + c.naming.spell(field.rpc)
}

// servesPath reports whether the last two segments of the package's path spell
// the service and rpc of field.
func (c *checker) servesPath(field rpcField) bool {
	return strings.HasSuffix(c.pass.Pkg.Path(), "/"+c.path(field))
}

// fieldsOf returns the <Rpc>Func fields of every generated service struct pkg
// declares, or none when pkg is not a generated connect package.
func (c *checker) fieldsOf(pkg *types.Package) []rpcField {
	if pkg == nil || !strings.HasSuffix(lastPathSegment(pkg.Path()), "connect") {
		return nil
	}
	if fields, ok := c.services[pkg]; ok {
		return fields
	}

	var fields []rpcField
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}
		named, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}
		fields = append(fields, serviceStructFields(named)...)
	}

	c.services[pkg] = fields
	return fields
}

// fieldFor reports the generated handler field t is the type of. A named type
// is looked up in the package that declares it. An unnamed func type, which is
// what a Handler declared as a plain func has, is compared against every
// handler type of the connect packages this package imports, and counts only
// when exactly one matches: two rpcs with the same signature leave no way to
// tell which of them the func is for.
func (c *checker) fieldFor(t types.Type) (rpcField, bool) {
	if named, ok := t.(*types.Named); ok {
		for _, field := range c.fieldsOf(named.Obj().Pkg()) {
			if types.Identical(field.typ, t) {
				return field, true
			}
		}
		return rpcField{}, false
	}

	if _, ok := t.(*types.Signature); !ok {
		return rpcField{}, false
	}

	var matches []rpcField
	for _, pkg := range c.pass.Pkg.Imports() {
		for _, field := range c.fieldsOf(pkg) {
			if types.Identical(field.typ.Underlying(), t) {
				matches = append(matches, field)
			}
		}
	}
	if len(matches) != 1 {
		return rpcField{}, false
	}
	return matches[0], true
}

// resultField reports the generated handler field that the first result of
// sig is the type of.
func (c *checker) resultField(sig *types.Signature) (rpcField, bool) {
	if sig.Results().Len() == 0 {
		return rpcField{}, false
	}
	return c.fieldFor(sig.Results().At(0).Type())
}

// handlerDecl is a top-level declaration of an rpc package whose type says it
// is, or builds, the handler of a generated service field.
type handlerDecl struct {
	name  *ast.Ident
	field rpcField
	// builds is true when the declaration returns the handler rather than
	// being it.
	builds bool
}

// checkHandlerPackage covers a package that declares an rpc handler: it lives
// at <service>/<rpc>, is named for the rpc, and declares its handler in
// <package>.go as either Handler or NewHandler, never both. Any other exported
// name that is or builds a handler is reported, since the wiring only ever
// looks for those two.
func (c *checker) checkHandlerPackage() {
	files := c.nonTestFiles()
	if len(files) == 0 {
		return
	}

	decls := c.handlerDecls(files)
	if !c.isRPCPackage(decls) {
		return
	}

	var handler, constructor *handlerDecl
	for i := range decls {
		decl := &decls[i]
		switch {
		case decl.name.Name == handlerName && !decl.builds:
			handler = decl
		case decl.name.Name == newHandlerName && decl.builds:
			constructor = decl
		case decl.builds:
			c.pass.Reportf(decl.name.Pos(), "%s builds the handler for rpc %s, so it must be named %s",
				decl.name.Name, decl.field.rpc, newHandlerName)
		default:
			c.pass.Reportf(decl.name.Pos(), "%s is the handler for rpc %s, so it must be named %s",
				decl.name.Name, decl.field.rpc, handlerName)
		}
	}

	declared := handler
	if declared == nil {
		declared = constructor
	}
	if declared == nil {
		declared = &decls[0]
	}
	field := declared.field

	if !c.servesPath(field) {
		c.pass.Reportf(files[0].Package, "package %s serves rpc %s of %s, so it must be at %s",
			c.pass.Pkg.Path(), field.rpc, field.service, c.path(field))
	}

	// The file is only worth naming once the package is: against a wrong
	// package name, <package>.go would point at a file that should not exist.
	wantName := c.naming.spell(field.rpc)
	nameMatches := c.pass.Pkg.Name() == wantName
	if !nameMatches {
		c.pass.Reportf(files[0].Package, "package %s serves rpc %s, so it must be named %s, got %s",
			c.pass.Pkg.Path(), field.rpc, wantName, c.pass.Pkg.Name())
	}

	// Both is the ambiguous case rather than the generous one: it leaves two
	// answers to what this rpc's handler is, and the wiring picks one, so the
	// package no longer says on its own which of them ships.
	if handler != nil && constructor != nil {
		c.pass.Reportf(constructor.name.Pos(), "package %s must declare %s or %s, not both",
			c.pass.Pkg.Path(), handlerName, newHandlerName)
		return
	}
	if handler == nil && constructor == nil || !nameMatches {
		return
	}

	wantFilename := c.pass.Pkg.Name() + ".go"
	gotFilename := filepath.Base(c.pass.Fset.Position(declared.name.Pos()).Filename)
	if gotFilename != wantFilename {
		c.pass.Reportf(declared.name.Pos(), "%s must be declared in %s, not in %s", declared.name.Name, wantFilename, gotFilename)
	}
}

// isRPCPackage reports whether decls make this an rpc package: one of them is
// named Handler or NewHandler, or the package sits at the path of the rpc one
// of them serves, in any spelling, so a misnamed handler is still found in a
// directory spelled the other way. Anything else that returns a handler, a wiring
// package handing a handler on for instance, is not an rpc package, and
// holding it to the rules for one would report the code that uses handlers
// rather than the code that declares them.
func (c *checker) isRPCPackage(decls []handlerDecl) bool {
	for _, decl := range decls {
		if decl.name.Name == handlerName || decl.name.Name == newHandlerName || nearPath(c.pass.Pkg.Path(), decl.field) {
			return true
		}
	}
	return false
}

// handlerDecls finds every exported top-level func or var that is, or whose
// first result is, one of the generated handler types. A func counts as the
// handler when its own signature is the handler type, and as building one
// when its first result is.
func (c *checker) handlerDecls(files []*ast.File) []handlerDecl {
	var decls []handlerDecl
	for _, f := range files {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Recv != nil || !decl.Name.IsExported() {
					continue
				}
				fn, ok := c.pass.TypesInfo.Defs[decl.Name].(*types.Func)
				if !ok {
					continue
				}
				sig := fn.Type().(*types.Signature)
				if field, ok := c.resultField(sig); ok {
					decls = append(decls, handlerDecl{name: decl.Name, field: field, builds: true})
				} else if field, ok := c.fieldFor(sig); ok {
					decls = append(decls, handlerDecl{name: decl.Name, field: field})
				}
			case *ast.GenDecl:
				if decl.Tok != token.VAR {
					continue
				}
				for _, spec := range decl.Specs {
					for _, name := range spec.(*ast.ValueSpec).Names {
						if !name.IsExported() {
							continue
						}
						obj, ok := c.pass.TypesInfo.Defs[name].(*types.Var)
						if !ok {
							continue
						}
						if field, ok := c.fieldFor(obj.Type()); ok {
							decls = append(decls, handlerDecl{name: name, field: field})
						}
					}
				}
			}
		}
	}
	return decls
}

// serviceStructFields returns the <Rpc>Func fields of named when it is a
// generated service struct: a struct type whose name ends with "Struct" and
// that has at least one <Rpc>Func field of func type, in declaration order.
func serviceStructFields(named *types.Named) []rpcField {
	structName := named.Obj().Name()
	if !strings.HasSuffix(structName, "Struct") {
		return nil
	}
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	var fields []rpcField
	for f := range structType.Fields() {
		if !strings.HasSuffix(f.Name(), "Func") || !isFuncType(f.Type()) {
			continue
		}
		fields = append(fields, rpcField{
			service: strings.TrimSuffix(structName, "Struct"),
			rpc:     strings.TrimSuffix(f.Name(), "Func"),
			typ:     f.Type(),
		})
	}
	return fields
}

// nearPath reports whether the last two segments of path name the service and
// rpc of field, with case and underscores ignored.
func nearPath(path string, field rpcField) bool {
	segments := strings.Split(path, "/")
	if len(segments) < 2 {
		return false
	}
	return normalize(segments[len(segments)-2]) == normalize(field.service) &&
		normalize(segments[len(segments)-1]) == normalize(field.rpc)
}

// normalize folds the spellings a service or rpc name is written in, so
// EchoStream, echostream and echo_stream compare equal.
func normalize(name string) string {
	return strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(name))
}

// snakeCase spells a Go identifier as snake_case words: EchoStream is
// echo_stream, PingV2Service is ping_v2_service, and an initialism stays one
// word, so GetHTTPStatus is get_http_status.
func snakeCase(name string) string {
	var b strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			startsWord := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if !unicode.IsUpper(runes[i-1]) || startsWord {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func isFuncType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Signature)
	return ok
}

func lastPathSegment(path string) string {
	if i := strings.LastIndexByte(path, '/'); i >= 0 {
		return path[i+1:]
	}
	return path
}

// nonTestFiles returns the package's files other than *_test.go, ordered by
// filename, so which file a report lands on does not depend on the driver's
// file ordering. None means this pass is an external test package, whose
// handler, if any, is checked in the pass of the package it tests.
func (c *checker) nonTestFiles() []*ast.File {
	var files []*ast.File
	for _, f := range c.pass.Files {
		if !strings.HasSuffix(c.pass.Fset.Position(f.Pos()).Filename, "_test.go") {
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return c.pass.Fset.Position(files[i].Pos()).Filename < c.pass.Fset.Position(files[j].Pos()).Filename
	})
	return files
}
