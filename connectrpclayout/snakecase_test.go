package connectrpclayout

import "testing"

func TestSnakeCase(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"Echo":              "echo",
		"EchoStream":        "echo_stream",
		"CheckoutV3Service": "checkout_v3_service",
		"GetHTTPStatus":     "get_http_status",
		"HTTPStatus":        "http_status",
		"GetURL":            "get_url",
	}

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := snakeCase(name); got != want {
				t.Fatalf("snakeCase(%q) = %q, want %q", name, got, want)
			}
		})
	}
}
