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

type Cache struct{}

func (Cache) Get(key string) error                              { return nil }
func (*Cache) GetContext(ctx context.Context, key string) error { return nil }

func newCache() Cache { return Cache{} }

type Inner struct{}

func (Inner) PutContext(ctx context.Context, key string) error { return nil }

type Outer struct{ Inner }

func (Outer) Put(key string) error { return nil }

type Wrapped struct{ Store }

func Map[T any](v T) T                             { return v }
func MapContext[U any](ctx context.Context, v U) U { return v }

type List[T any] struct{}

func (List[T]) Push(v T)                             {}
func (List[T]) PushContext(ctx context.Context, v T) {}

func receivers(c Cache, o Outer, w Wrapped, l List[int]) {
	_ = c.Get("a") // want `\(Cache\).Get drops the context; call GetContext`
	_ = newCache().Get("a")
	_ = o.Put("a")
	_ = w.Load("a") // want `\(Store\).Load drops the context; call LoadContext`
	_ = Map(1)      // want `^Map drops the context; call MapContext`
	l.Push(1)       // want `\(List\[T\]\).Push drops the context; call PushContext`
}
