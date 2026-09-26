package main

import "context"

func main() {
	ctx := context.Background()
	go func() {
		_ = context.TODO() // want `context.TODO starts a new context`
	}()
	_ = ctx
}

func helper() {
	_ = context.Background() // want `context.Background starts a new context`
}
