// Package appsemconv is a project's own home for the names semconv does not
// have. It spells them itself, which is the one place that is allowed to once
// it is listed in allowed-packages.
package appsemconv

import "go.opentelemetry.io/otel/attribute"

const TierKey = attribute.Key("app.user.tier")

const TierName = "app.user.tier"

const RequestCountName = "app.request.count"

func Tier(v string) attribute.KeyValue { return attribute.String("app.user.tier", v) }
