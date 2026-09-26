// Package attribute stands in for go.opentelemetry.io/otel/attribute, with
// only the shapes semconvkey looks at.
package attribute

type Key string

type KeyValue struct {
	Key   Key
	Value string
}

func (k Key) String(v string) KeyValue { return KeyValue{Key: k, Value: v} }

func (k Key) Bool(bool) KeyValue { return KeyValue{Key: k} }

func (k Key) Int64(int64) KeyValue { return KeyValue{Key: k} }

func (k Key) Defined() bool { return k != "" }

func String(k, v string) KeyValue { return KeyValue{Key: Key(k), Value: v} }

func Bool(k string, _ bool) KeyValue { return KeyValue{Key: Key(k)} }

func Int64(k string, _ int64) KeyValue { return KeyValue{Key: Key(k)} }
