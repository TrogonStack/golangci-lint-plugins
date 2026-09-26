# staticerr

Reports any `errors.New` call outside a package-level `var`, exported or not.
An error built where it is returned is a new value on every call, so
`errors.Is` can never match it. `_test.go` files are skipped.
