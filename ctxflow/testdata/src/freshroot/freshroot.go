package freshroot

import (
	"context"
	stdctx "context"
	"time"
)

var shared = context.Background() // want `context.Background starts a new context`

func run() {
	_ = context.Background()     // want `context.Background starts a new context`
	_ = context.TODO()           // want `context.TODO starts a new context`
	_ = stdctx.Background()      // want `context.Background starts a new context`
	newCtx := context.Background // want `context.Background starts a new context`
	_ = newCtx
}

func detach(ctx context.Context) {
	flush, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
	defer cancel()
	_ = flush
}

func init() {
	_ = context.Background()
}
