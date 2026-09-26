// Package semconv stands in for go.opentelemetry.io/otel/semconv, with a key,
// a KeyValue constructor and metric names in each of the shapes upstream
// declares them.
package semconv

import "go.opentelemetry.io/otel/attribute"

const MessagingSystemKey = attribute.Key("messaging.system")

func ServerAddress(v string) attribute.KeyValue {
	return attribute.KeyValue{Key: "server.address", Value: v}
}

const ServerActiveRequestsName = "http.server.active_requests"

type ServerRequestDuration struct{}

func (ServerRequestDuration) Name() string { return "http.server.request.duration" }
