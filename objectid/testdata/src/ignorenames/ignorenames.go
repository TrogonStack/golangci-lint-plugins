// Package ignorenames is a fixture for ignore-names, which matches a name
// exactly and case-sensitively.
package ignorenames

func Open(projectID, databaseID string, orderID string) {} // want `orderID is a raw string; use a typed ID`

type Config struct {
	ProjectID  string // want `ProjectID is a raw string; use a typed ID`
	databaseID string
}
