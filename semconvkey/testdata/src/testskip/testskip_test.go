package testskip

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestTier(t *testing.T) {
	if got := attribute.String("app.user.tier", "gold"); got.Key != "app.user.tier" {
		t.Fatal(got)
	}
}
