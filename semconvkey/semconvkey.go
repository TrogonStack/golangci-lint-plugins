// Package semconvkey checks that an OpenTelemetry attribute key or metric
// instrument name comes from a package allowed to declare one, rather than
// being spelled by hand at the call site. go.opentelemetry.io and everything
// under it is always allowed, so upstream semconv works with no
// configuration; a project that declares its own names lists the packages
// holding them in allowed-packages, each named to end in semconv.
//
// The reason is that an attribute key or a metric name is a contract in the
// same way a protobuf field is: a dashboard, an alert or a query built
// against it breaks silently the day it is respelled. Declaring each name
// once, in a package whose job is to hold them, is what keeps it from
// drifting; a literal string or a local const at a call site is a second,
// unreviewed place the same name can change.
//
// It reports, everywhere but an allowed package:
//
//   - an attribute constructor that takes the key as a raw string, such as
//     attribute.String("app.user.tier", v), whatever the string is;
//   - a method of attribute.Key that builds a KeyValue, such as
//     key.String(v), unless the key comes from an allowed package;
//   - a conversion to attribute.Key, unless its operand comes from an allowed
//     package;
//   - an attribute.KeyValue composite literal, unless its Key does;
//   - a metric.Meter method that creates an instrument, such as
//     meter.Int64Counter(name), unless the name does, or whatever the name
//     is when generated-instruments is set;
//   - a composite literal of a struct type a package listed in
//     allowed-packages declares, when every field of that type is
//     unexported, such as appsemconv.TierAttr{}.
//
// A literal of such a type written outside its package can only ever be the
// zero value, since there is no field it could set: an attribute with its
// key and an empty value, or an instrument that records nothing. The
// package's own constructor is the only way to get one that means something,
// and a parameter typed as one is only required if the zero value cannot be
// passed in its place. go.opentelemetry.io is left out of this rule, since
// it declares many such types whose zero value is fine to build.
//
// A value comes from an allowed package when it is a constant, variable,
// field or function result that package declares, named directly at the use.
// A local copy of one does not count, even an unchanged one: the rule is
// about where a name is spelled, and a local is a place it could be
// respelled. *_test.go files and generated files are skipped: a test asserts
// the wire name as a string on purpose, and generated code is never
// hand-edited to begin with.
package semconvkey

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Doc is the analyzer's one-line description, shown by go vet -help and
// golangci-lint's linter listing.
const Doc = "checks that OpenTelemetry attribute keys and metric instrument names come from semconv or an allowed package, not a literal at the call site"

const (
	attributePackage = "go.opentelemetry.io/otel/attribute"
	metricPackage    = "go.opentelemetry.io/otel/metric"
)

// upstream is always allowed, since it is where semconv and the
// OpenTelemetry API that builds on it declare their own names.
const upstream PackagePrefix = "go.opentelemetry.io"

const (
	keyMessage         = "an OpenTelemetry attribute key must come from go.opentelemetry.io/otel/semconv or a package listed in allowed-packages, not be spelled at the call site"
	rawKeyMessage      = "attribute.%s takes the key as a raw string; call the method of a Key declared in go.opentelemetry.io/otel/semconv or a package listed in allowed-packages instead"
	metricNameMessage  = "an OpenTelemetry metric name must come from go.opentelemetry.io/otel/semconv or a package listed in allowed-packages, not be spelled at the call site"
	instrumentMessage  = "an OpenTelemetry metric instrument must be created by a package listed in allowed-packages, so its name, unit and required attributes come with it"
	zeroValueMessage   = "%s has no exported fields, so a composite literal of it can only be its zero value; build it with the constructor %s declares"
	generatedHeaderRaw = `^// Code generated .* DO NOT EDIT\.$`
)

// instrumentMethods is every metric.Meter method that mints an instrument
// identified by name, synchronous and observable alike. RegisterCallback is
// missing on purpose: it takes instruments, not a name.
var instrumentMethods = map[string]bool{
	"Int64Counter":                   true,
	"Int64UpDownCounter":             true,
	"Int64Histogram":                 true,
	"Int64Gauge":                     true,
	"Int64ObservableCounter":         true,
	"Int64ObservableUpDownCounter":   true,
	"Int64ObservableGauge":           true,
	"Float64Counter":                 true,
	"Float64UpDownCounter":           true,
	"Float64Histogram":               true,
	"Float64Gauge":                   true,
	"Float64ObservableCounter":       true,
	"Float64ObservableUpDownCounter": true,
	"Float64ObservableGauge":         true,
}

// generatedHeader is the convention https://go.dev/s/generatedcode defines.
var generatedHeader = regexp.MustCompile(generatedHeaderRaw)

var (
	ErrEmptyPackagePrefix      = errors.New("allowed-packages entry is empty")
	ErrPackagePrefixNotSemconv = errors.New("allowed-packages entry must be a package whose name ends in semconv")
)

// semconvSuffix is what an allowed package's name must end in, so a reader
// can tell a package that declares names apart from one that only uses them.
const semconvSuffix = "semconv"

// PackagePrefix is an import path that, with everything under it, is allowed
// to declare attribute keys and metric names.
type PackagePrefix string

// covers reports whether path is the prefix itself or a package under it.
func (p PackagePrefix) covers(path string) bool {
	return path == string(p) || strings.HasPrefix(path, string(p)+"/")
}

func (p PackagePrefix) namedForSemconv() bool {
	return strings.HasSuffix(path.Base(string(p)), semconvSuffix)
}

// Settings is what golangci-lint's custom linter settings decode into.
type Settings struct {
	AllowedPackages []PackagePrefix `json:"allowed-packages"`

	// GeneratedInstruments reports every instrument created outside an
	// allowed package, whatever its name, for a project whose allowed
	// packages build its instruments rather than only naming them.
	GeneratedInstruments bool `json:"generated-instruments"`
}

// Analyzer is semconvkey with the default settings, allowing only
// go.opentelemetry.io, for any analysis driver: golangci-lint through
// semconvkey/plugin, or a singlechecker binary.
var Analyzer = newAnalyzer(checker{allowed: allowed{upstream}})

// New is semconvkey configured by settings.
func New(settings Settings) (*analysis.Analyzer, error) {
	var configured allowed
	for i, prefix := range settings.AllowedPackages {
		if prefix == "" {
			return nil, fmt.Errorf("%w: entry %d", ErrEmptyPackagePrefix, i)
		}
		if !prefix.namedForSemconv() {
			return nil, fmt.Errorf("%w: %q", ErrPackagePrefixNotSemconv, prefix)
		}
		configured = append(configured, prefix)
	}
	return newAnalyzer(checker{
		allowed:              append(allowed{upstream}, configured...),
		configured:           configured,
		generatedInstruments: settings.GeneratedInstruments,
	}), nil
}

type allowed []PackagePrefix

func (a allowed) covers(pkg *types.Package) bool {
	if pkg == nil {
		return false
	}
	for _, prefix := range a {
		if prefix.covers(pkg.Path()) {
			return true
		}
	}
	return false
}

func newAnalyzer(config checker) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "semconvkey",
		Doc:  Doc,
		Run: func(pass *analysis.Pass) (any, error) {
			c := config
			c.pass = pass
			c.run()
			return nil, nil
		},
	}
}

type checker struct {
	pass                 *analysis.Pass
	allowed              allowed
	configured           allowed
	generatedInstruments bool
}

func (c checker) run() {
	if c.allowed.covers(c.pass.Pkg) {
		return
	}

	for _, file := range c.pass.Files {
		filename := c.pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") || isGenerated(file) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				c.checkCall(node)
			case *ast.CompositeLit:
				c.checkCompositeLit(node)
			}
			return true
		})
	}
}

func (c checker) checkCall(call *ast.CallExpr) {
	if c.isKeyConversion(call) {
		if len(call.Args) == 1 && !c.fromAllowed(call.Args[0]) {
			c.pass.Reportf(call.Pos(), "%s", keyMessage)
		}
		return
	}

	sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	if !ok {
		return
	}
	fn, ok := c.pass.TypesInfo.Uses[sel.Sel].(*types.Func)
	if !ok || fn.Pkg() == nil {
		return
	}
	sig := fn.Type().(*types.Signature)

	switch fn.Pkg().Path() {
	case attributePackage:
		if !returnsKeyValue(sig) {
			return
		}
		if sig.Recv() == nil {
			if takesRawKey(sig) {
				c.pass.Reportf(call.Pos(), rawKeyMessage, fn.Name())
			}
			return
		}
		if !c.isKeyConversion(sel.X) && !c.fromAllowed(sel.X) {
			c.pass.Reportf(call.Pos(), "%s", keyMessage)
		}
	case metricPackage:
		if !instrumentMethods[fn.Name()] {
			return
		}
		if c.generatedInstruments {
			c.pass.Reportf(call.Pos(), "%s", instrumentMessage)
			return
		}
		if len(call.Args) > 0 && !c.fromAllowed(call.Args[0]) {
			c.pass.Reportf(call.Pos(), "%s", metricNameMessage)
		}
	}
}

// isKeyConversion reports whether expr is a conversion to attribute.Key,
// which checkCall judges on its own, so a method called on one is not
// reported a second time for the same key.
func (c checker) isKeyConversion(expr ast.Expr) bool {
	call, ok := ast.Unparen(expr).(*ast.CallExpr)
	return ok && isAttributeType(c.typeOfCallee(call), "Key")
}

// typeOfCallee is the type call.Fun names when the call is a conversion, or
// nil when it is a call of a function.
func (c checker) typeOfCallee(call *ast.CallExpr) types.Type {
	tv, ok := c.pass.TypesInfo.Types[call.Fun]
	if !ok || !tv.IsType() {
		return nil
	}
	return tv.Type
}

func (c checker) checkCompositeLit(lit *ast.CompositeLit) {
	t := c.pass.TypesInfo.TypeOf(lit)

	if named, ok := c.opaqueConfiguredType(t); ok {
		c.pass.Reportf(lit.Pos(), zeroValueMessage, named.Obj().Name(), named.Obj().Pkg().Name())
		return
	}

	if !isAttributeType(t, "KeyValue") {
		return
	}

	key := compositeKey(lit)
	if key == nil || !c.fromAllowed(key) {
		c.pass.Reportf(lit.Pos(), "%s", keyMessage)
	}
}

// opaqueConfiguredType reports whether t is a struct type declared in a
// configured allowed package with at least one field and none exported. An
// empty struct is left alone: its zero value is the only value it has.
func (c checker) opaqueConfiguredType(t types.Type) (*types.Named, bool) {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || !c.configured.covers(named.Obj().Pkg()) {
		return nil, false
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok || st.NumFields() == 0 {
		return nil, false
	}

	for field := range st.Fields() {
		if field.Exported() {
			return nil, false
		}
	}

	return named, true
}

// compositeKey is the expression a KeyValue literal sets its Key to, keyed or
// positional, or nil when it leaves Key unset.
func compositeKey(lit *ast.CompositeLit) ast.Expr {
	for i, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			if i == 0 {
				return elt
			}
			continue
		}
		if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "Key" {
			return kv.Value
		}
	}
	return nil
}

// fromAllowed reports whether expr names something an allowed package
// declares: one of its constants, variables or fields, or the result of
// calling one of its functions or methods, or a conversion of one. Anything
// else, a literal, a local, a concatenation, is a name spelled here.
func (c checker) fromAllowed(expr ast.Expr) bool {
	switch e := ast.Unparen(expr).(type) {
	case *ast.Ident:
		return c.declaredInAllowed(e)
	case *ast.SelectorExpr:
		return c.declaredInAllowed(e.Sel)
	case *ast.CallExpr:
		if c.typeOfCallee(e) != nil {
			return len(e.Args) == 1 && c.fromAllowed(e.Args[0])
		}
		switch fun := ast.Unparen(e.Fun).(type) {
		case *ast.Ident:
			return c.declaredInAllowed(fun)
		case *ast.SelectorExpr:
			return c.declaredInAllowed(fun.Sel)
		}
	}
	return false
}

func (c checker) declaredInAllowed(ident *ast.Ident) bool {
	obj := c.pass.TypesInfo.Uses[ident]
	return obj != nil && c.allowed.covers(obj.Pkg())
}

// returnsKeyValue reports whether sig returns exactly one attribute.KeyValue.
// Matching by signature rather than by name catches every constructor and
// every Key method, including ones upstream adds later.
func returnsKeyValue(sig *types.Signature) bool {
	return sig.Results().Len() == 1 && isAttributeType(sig.Results().At(0).Type(), "KeyValue")
}

// takesRawKey reports whether sig's first parameter is the key as a plain
// string, which is what attribute.String, attribute.Int and the rest take.
func takesRawKey(sig *types.Signature) bool {
	if sig.Params().Len() == 0 {
		return false
	}
	basic, ok := sig.Params().At(0).Type().(*types.Basic)
	return ok && basic.Kind() == types.String
}

func isAttributeType(t types.Type, name string) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == attributePackage && obj.Name() == name
}

// isGenerated matches the comment https://go.dev/s/generatedcode asks a
// generator to write, in the comments above the package clause.
func isGenerated(file *ast.File) bool {
	for _, group := range file.Comments {
		if group.Pos() > file.Package {
			break
		}
		for _, comment := range group.List {
			if generatedHeader.MatchString(comment.Text) {
				return true
			}
		}
	}
	return false
}
