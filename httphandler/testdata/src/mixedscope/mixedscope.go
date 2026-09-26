// Package mixedscope is a fixture for confirming that httphandler ignores a
// NewHandler outside its scope when judging whether a package declares both
// forms.
package mixedscope

import (
	"io"
	"net/http"
)

type AuditLogHandler struct{}

func Handler(w http.ResponseWriter, r *http.Request) {}

func NewHandler(w io.Writer) *AuditLogHandler { return nil }
