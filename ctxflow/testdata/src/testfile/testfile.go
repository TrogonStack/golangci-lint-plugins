package testfile

import "context"

func main() {
	_ = context.Background() // want `context.Background starts a new context`
}
