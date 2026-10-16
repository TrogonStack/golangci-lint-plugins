// Package ignorenamesoff is ignorenames with no ignore-names set, so the
// same names are reported.
package ignorenamesoff

func Open(projectID, databaseID string) {} // want `projectID is a raw string; use a typed ID` `databaseID is a raw string; use a typed ID`
