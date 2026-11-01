package attributedotcapture

import . "go.opentelemetry.io/otel/attribute"

func Capture() {
	create := String // want `takes the key as a raw string`
	_ = create("app.raw", "value")
}
