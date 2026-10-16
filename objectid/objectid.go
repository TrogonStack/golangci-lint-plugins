// Package objectid checks that an identifier is declared with a type of its
// own, such as OrderID, rather than as a raw string, integer or uuid.UUID.
//
// It reports a function parameter, a function result, an interface method's
// parameter or result, and a struct field when:
//
//   - its name matches name-patterns, by default a name ending in ID, Id,
//     IDs or Ids, or the bare name id or ids; and
//   - its type, after any pointers, slices and arrays are taken off, is the
//     predeclared string type, a predeclared integer type, or
//     github.com/google/uuid.UUID.
//
// A named type such as type OrderID string is not reported: declaring one is
// the fix. An alias of string is reported, since it is string by another
// name.
//
// The reason is that two IDs of the same raw type are interchangeable to the
// compiler. A function taking (orderID, userID string) compiles just as well
// when a caller swaps them, and the bug shows up as a row that is not found,
// or worse, as the wrong row. A type per ID makes the swap a compile error,
// and gives the one place an ID is parsed and validated a home.
//
// With untyped-constants set, it also reports an untyped constant that
// becomes a named ID type outside a const declaration, such as
// OrderID("abc") or a literal "abc" passed where an OrderID is expected. A
// named type stops a raw string from being passed by mistake, but not a
// literal, which converts to it implicitly. The zero value, "" or 0, is left
// alone, since comparing against it is how an unset ID is detected.
//
// Skipped:
//
//   - generated files and *_test.go files;
//   - packages under an exclude-packages entry;
//   - names listed in ignore-names, matched exactly;
//   - function literals, whose signature is usually imposed by whatever the
//     literal is passed to;
//   - a method whose receiver implements an interface declared in a package
//     this one imports, directly or not, when that interface declares a
//     method of the same name and an identical signature. The signature is
//     the interface's, so the interface is where it is reported. An
//     interface declared in a package that imports this one is not seen.
package objectid

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that parameters, results and struct fields named like an ID have a typed ID instead of a raw string, integer or uuid.UUID"

const (
	rawMessage      = "%s is a raw %s; use a typed ID%s"
	untypedMessage  = "untyped constant converted to %s; build the ID with its constructor or parser instead of a literal%s"
	uuidPackagePath = "github.com/google/uuid"
	uuidTypeName    = "UUID"
)

// defaultNamePatterns are the suffixes Go naming uses for an identifier
// and a list of them, in initialisms and in mixed case.
var defaultNamePatterns = []string{"ID", "Id", "IDs", "Ids"}

var (
	ErrEmptyNamePattern   = errors.New("name-patterns entry is empty")
	ErrEmptyIgnoreName    = errors.New("ignore-names entry is empty")
	ErrEmptyPackagePrefix = errors.New("exclude-packages entry is empty")
)

// PackagePrefix is an import path that, with everything under it, is not
// checked.
type PackagePrefix string

// covers reports whether path is the prefix itself or a package under it.
func (p PackagePrefix) covers(path string) bool {
	return path == string(p) || strings.HasPrefix(path, string(p)+"/")
}

// Settings is what golangci-lint's custom linter settings decode into.
type Settings struct {
	// NamePatterns are the name suffixes that mark an ID. A name equal to a
	// suffix in lower case, such as id, matches too. Empty means the
	// defaults: ID, Id, IDs and Ids.
	NamePatterns []string `json:"name-patterns"`

	// IgnoreNames are exact, case-sensitive names that are never reported,
	// for a name that matches name-patterns but is not a domain ID, such as
	// a cloud project or database name.
	IgnoreNames []string `json:"ignore-names"`

	// Suggestion is appended to the message in parentheses, to name where
	// the project's typed IDs live.
	Suggestion string `json:"suggestion"`

	// ExcludePackages are import paths that, with every package under them,
	// are not checked.
	ExcludePackages []PackagePrefix `json:"exclude-packages"`

	// UntypedConstants also reports an untyped constant converted to a named
	// ID type outside a test or a const declaration.
	UntypedConstants bool `json:"untyped-constants"`
}

// Analyzer is objectid with the default settings, for any analysis driver:
// golangci-lint through objectid/plugin, or a singlechecker binary.
var Analyzer = newAnalyzer(config{patterns: defaultNamePatterns})

// New is objectid configured by settings.
func New(settings Settings) (*analysis.Analyzer, error) {
	names := defaultNamePatterns
	if len(settings.NamePatterns) > 0 {
		names = settings.NamePatterns
	}
	for i, name := range names {
		if name == "" {
			return nil, fmt.Errorf("%w: entry %d", ErrEmptyNamePattern, i)
		}
	}
	ignore := make(map[string]bool, len(settings.IgnoreNames))
	for i, name := range settings.IgnoreNames {
		if name == "" {
			return nil, fmt.Errorf("%w: entry %d", ErrEmptyIgnoreName, i)
		}
		ignore[name] = true
	}
	for i, prefix := range settings.ExcludePackages {
		if prefix == "" {
			return nil, fmt.Errorf("%w: entry %d", ErrEmptyPackagePrefix, i)
		}
	}
	var suffix string
	if settings.Suggestion != "" {
		suffix = " (" + settings.Suggestion + ")"
	}
	return newAnalyzer(config{
		patterns:         names,
		ignore:           ignore,
		suffix:           suffix,
		exclude:          settings.ExcludePackages,
		untypedConstants: settings.UntypedConstants,
	}), nil
}

// patterns are name suffixes; a name equal to one in lower case matches too.
type patterns []string

// match is case-sensitive, so that valid and Grid, which end in id, do not
// read as IDs.
func (p patterns) match(name string) bool {
	for _, s := range p {
		if name == strings.ToLower(s) || strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

type config struct {
	patterns         patterns
	ignore           map[string]bool
	suffix           string // the suggestion, in parentheses, for every message
	exclude          []PackagePrefix
	untypedConstants bool
}

func newAnalyzer(cfg config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "objectid",
		Doc:  Doc,
		Run: func(pass *analysis.Pass) (any, error) {
			c := &checker{config: cfg, pass: pass}
			c.run()
			return nil, nil
		},
	}
}

type checker struct {
	config
	pass *analysis.Pass

	// interfaces indexes, by method name, the interfaces of every package
	// this one imports. It is built on first use, since most packages have
	// no method with a raw ID to look up.
	interfaces map[string][]*types.Interface
}

func (c *checker) run() {
	for _, prefix := range c.exclude {
		if prefix.covers(c.pass.Pkg.Path()) {
			return
		}
	}

	for _, file := range c.pass.Files {
		filename := c.pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") || ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, c.visit)
	}
}

// visit checks n and reports whether to descend into it.
func (c *checker) visit(n ast.Node) bool {
	switch node := n.(type) {
	case *ast.FuncDecl:
		found := append(c.rawFields(node.Type.Params), c.rawFields(node.Type.Results)...)
		if len(found) > 0 && !c.imposed(node) {
			c.report(found)
		}
	case *ast.FuncLit:
		// Skip the literal's own signature but keep going into its body,
		// which can declare struct types of its own.
		ast.Inspect(node.Body, c.visit)
		return false
	case *ast.StructType:
		c.report(c.rawFields(node.Fields))
	case *ast.InterfaceType:
		for _, method := range node.Methods.List {
			if ft, ok := method.Type.(*ast.FuncType); ok {
				c.report(c.rawFields(ft.Params))
				c.report(c.rawFields(ft.Results))
			}
		}
	case *ast.GenDecl:
		if node.Tok == token.CONST {
			// A typed constant is how a well-known ID is declared, so its
			// value is not a stray literal.
			return false
		}
	case ast.Expr:
		if c.untypedConstants && c.checkUntyped(node) {
			return false
		}
	}
	return true
}

// rawField is a field name that would be reported, with its type spelled
// for the message.
type rawField struct {
	name *ast.Ident
	raw  string
}

func (c *checker) rawFields(fields *ast.FieldList) []rawField {
	if fields == nil {
		return nil
	}
	var found []rawField
	for _, field := range fields.List {
		// An embedded field has no Names, and is skipped on purpose: a struct
		// embedding uuid.UUID is how a typed ID is declared.
		for _, name := range field.Names {
			if name.Name == "_" || c.ignore[name.Name] || !c.patterns.match(name.Name) {
				continue
			}
			obj := c.pass.TypesInfo.Defs[name]
			if obj == nil {
				continue
			}
			raw, ok := rawID(obj.Type())
			if !ok {
				continue
			}
			found = append(found, rawField{name: name, raw: raw})
		}
	}
	return found
}

func (c *checker) report(found []rawField) {
	for _, f := range found {
		c.pass.Reportf(f.name.Pos(), rawMessage, f.name.Name, f.raw, c.suffix)
	}
}

// rawID reports whether t, once its pointers, slices and arrays are taken
// off, is a predeclared string or integer type, or uuid.UUID, and spells it
// with aliases resolved, so a message names what the type really is. A named
// type with one of those underneath is not raw: it is what the rule asks for.
func rawID(t types.Type) (string, bool) {
	var prefix strings.Builder
	for {
		switch u := types.Unalias(t).(type) {
		case *types.Pointer:
			prefix.WriteString("*")
			t = u.Elem()
		case *types.Slice:
			prefix.WriteString("[]")
			t = u.Elem()
		case *types.Array:
			fmt.Fprintf(&prefix, "[%d]", u.Len())
			t = u.Elem()
		case *types.Basic:
			if u.Info()&(types.IsString|types.IsInteger) == 0 {
				return "", false
			}
			return prefix.String() + u.Name(), true
		case *types.Named:
			obj := u.Obj()
			if obj.Pkg() == nil || obj.Pkg().Path() != uuidPackagePath || obj.Name() != uuidTypeName {
				return "", false
			}
			return prefix.String() + "uuid." + uuidTypeName, true
		default:
			return "", false
		}
	}
}

// imposed reports whether fn is a method whose signature an interface from
// an imported package dictates: the interface has a method of the same name
// and an identical signature, and the receiver's type implements it.
func (c *checker) imposed(fn *ast.FuncDecl) bool {
	if fn.Recv == nil {
		return false
	}
	obj, ok := c.pass.TypesInfo.Defs[fn.Name].(*types.Func)
	if !ok {
		return false
	}
	recv := obj.Signature().Recv().Type()
	if ptr, ok := types.Unalias(recv).(*types.Pointer); ok {
		recv = ptr.Elem()
	}
	// The pointer's method set holds every method, so it implements an
	// interface whenever either the value or the pointer does.
	recvPtr := types.NewPointer(recv)

	if c.interfaces == nil {
		c.interfaces = indexInterfaces(c.pass.Pkg)
	}
	// Implements needs the receiver's method to match the interface's, so
	// an interface found by name is one whose signature fn shares.
	for _, iface := range c.interfaces[fn.Name.Name] {
		if types.Implements(recvPtr, iface) {
			return true
		}
	}
	return false
}

// indexInterfaces collects the non-generic interface types declared at
// package scope in every package pkg imports, directly or not, by the names
// of their methods.
func indexInterfaces(pkg *types.Package) map[string][]*types.Interface {
	index := map[string][]*types.Interface{}
	// An import cycle is a compile error, so pkg is never its own dependency.
	for _, p := range typeutil.Dependencies(pkg.Imports()...) {
		scope := p.Scope()
		for _, name := range scope.Names() {
			tn, ok := scope.Lookup(name).(*types.TypeName)
			if !ok || tn.IsAlias() {
				continue
			}
			named, ok := tn.Type().(*types.Named)
			if !ok || named.TypeParams().Len() > 0 {
				continue
			}
			iface, ok := named.Underlying().(*types.Interface)
			if !ok {
				continue
			}
			for m := range iface.Methods() {
				index[m.Name()] = append(index[m.Name()], iface)
			}
		}
	}
	return index
}

// checkUntyped reports expr when it is an untyped constant that becomes a
// named ID type, and reports whether expr was an untyped constant at all, so
// visit need not descend into it.
func (c *checker) checkUntyped(expr ast.Expr) bool {
	tv, ok := c.pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil || !c.isUntypedConstant(expr) {
		return false
	}
	named, ok := types.Unalias(tv.Type).(*types.Named)
	if ok && c.patterns.match(named.Obj().Name()) && !isZero(tv.Value) {
		c.pass.Reportf(expr.Pos(), untypedMessage, types.TypeString(named, types.RelativeTo(c.pass.Pkg)), c.suffix)
	}
	return true
}

// isUntypedConstant reports whether expr, a constant, is built only from
// literals and untyped constants. The type checker records such an
// expression with the type it was converted to, so its own type cannot
// tell; how it is spelled can. Every operand of a constant is a constant, so
// only the spelling is checked here.
func (c *checker) isUntypedConstant(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return c.untypedConstIdent(e)
	case *ast.SelectorExpr:
		return c.untypedConstIdent(e.Sel)
	case *ast.ParenExpr:
		return c.isUntypedConstant(e.X)
	case *ast.UnaryExpr:
		return c.isUntypedConstant(e.X)
	case *ast.BinaryExpr:
		return c.isUntypedConstant(e.X) && c.isUntypedConstant(e.Y)
	}
	return false
}

func (c *checker) untypedConstIdent(ident *ast.Ident) bool {
	obj, ok := c.pass.TypesInfo.Uses[ident].(*types.Const)
	if !ok {
		return false
	}
	basic, ok := obj.Type().(*types.Basic)
	return ok && basic.Info()&types.IsUntyped != 0
}

func isZero(v constant.Value) bool {
	switch v.Kind() {
	case constant.String:
		return constant.StringVal(v) == ""
	case constant.Int, constant.Float, constant.Complex:
		return constant.Sign(v) == 0
	}
	return false
}
