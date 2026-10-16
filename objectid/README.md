# objectid

Checks that an ID has a type of its own instead of a raw string, integer or
`uuid.UUID`. It reports a function parameter, a function result, an
interface method's parameter or result, and a struct field when:

- its name ends in one of `name-patterns`, by default `ID`, `Id`, `IDs` or
  `Ids`, or is one of them in lower case, such as `id` or `ids`. The match is
  case-sensitive, so `valid` and `Grid` are not IDs;
- its type, after pointers, slices and arrays are taken off, is `string`, an
  integer type, or `github.com/google/uuid.UUID`. An alias of one counts.

A named type, such as `type OrderID string`, is not reported: declaring one
is the fix.

Two IDs of the same raw type are interchangeable to the compiler, so a caller
that swaps them compiles, and the bug shows up as a missing or wrong row. A
type per ID makes the swap a compile error.

```go
type Refunds interface {
	Refund(orderID, userID string) error          // reported twice
	Reverse(orderID OrderID, userID UserID) error // ok
}

type Order struct {
	ID       uuid.UUID // reported
	BuyerIDs []string  // reported
	SellerID SellerID  // ok
}
```

With `untyped-constants`, it also reports an untyped constant that becomes a
named ID type, such as `OrderID("abc")`, or `"abc"` passed where an `OrderID`
is expected. A named type stops a raw string from being passed by mistake,
but not a literal, which converts to it implicitly. A `const` declaration is
not reported, since that is how a well-known ID is declared, and neither is
the zero value, `""` or `0`, since comparing against it is how an unset ID is
detected.

Skipped:

- `*_test.go` and generated files;
- packages under an `exclude-packages` entry;
- names listed in `ignore-names`, matched exactly and case-sensitively;
- function literals, since their signature is usually imposed by what the
  literal is passed to;
- a method whose receiver implements an interface from an imported package,
  directly or not, when that interface has a method of the same name and an
  identical signature. The interface owns the signature, so it is reported
  where the interface is declared. An interface in a package that imports
  this one is not seen, so a method that implements one is still reported.

## Settings

| Setting             | Values                                                                                          | Default                  |
| ------------------- | ----------------------------------------------------------------------------------------------- | ------------------------ |
| `name-patterns`     | Name suffixes that mark an ID. A name equal to a suffix in lower case matches too               | `ID`, `Id`, `IDs`, `Ids` |
| `ignore-names`      | Exact, case-sensitive names never reported, such as `projectID` for a cloud project name         | none                     |
| `suggestion`        | Text added to the message in parentheses, such as the package that holds the project's typed IDs | none                     |
| `exclude-packages`  | Import paths that, with every package under them, are not checked                               | none                     |
| `untyped-constants` | `true` to also report an untyped constant converted to a named ID type                          | `false`                  |

An empty `name-patterns`, `ignore-names` or `exclude-packages` entry fails plugin load.

```yaml
linters:
  settings:
    custom:
      objectid:
        type: module
        settings:
          ignore-names:
            - projectID
            - databaseID
          suggestion: internal/ids
          exclude-packages:
            - example.com/app/gen
          untyped-constants: true
```
