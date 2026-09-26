package testfile

import (
	"context"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = context.Background()
	os.Exit(m.Run())
}

func TestSomething(t *testing.T) {
	_ = context.Background() // want `context.Background starts a new context`
	_ = t.Context()
}
