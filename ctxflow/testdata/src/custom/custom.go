package custom

import "context"

type Store struct{}

func (Store) Load(key string) error                             { return nil }
func (Store) LoadContext(ctx context.Context, key string) error { return nil }

func (Store) Save(key string) error                          { return nil }
func (Store) SaveContext(ctx context.Context, key int) error { return nil }

type Loader interface {
	Fetch(key string) ([]byte, error)
	FetchWithContext(ctx context.Context, key string) ([]byte, error)
}

func Send(msg string)                             {}
func SendContext(ctx context.Context, msg string) {}

func Close()                                 {}
func CloseContext(ctx context.Context) error { return nil }

func run(ctx context.Context, s Store, l Loader) {
	_ = s.Load("a")     // want `\(Store\).Load drops the context; call LoadContext`
	_, _ = l.Fetch("a") // want `\(Loader\).Fetch drops the context; call FetchWithContext`
	Send("hi")          // want `^Send drops the context; call SendContext`

	_ = s.Save("a")
	Close()
	_ = s.LoadContext(ctx, "a")
	SendContext(ctx, "hi")
}
