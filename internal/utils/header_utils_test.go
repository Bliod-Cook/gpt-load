package utils

import (
	"net/http"
	"testing"
)

func TestRemoveUpstreamPrivacyHeaders_RemovesCaseInsensitive(t *testing.T) {
	headers := http.Header{
		"x-real-ip":          {"1.2.3.4"},
		"X-Original-URL":     {"https://example.com/v1/chat/completions"},
		"X-Original-Method":  {"POST"},
		"X-Forwarded-URI":    {"/v1/chat/completions"},
		"X-Forwarded-Scheme": {"https"},
		"X-Forwarded-Method": {"POST"},
		"Content-Type":       {"application/json"},
		"X-Custom":           {"1"},
	}

	RemoveUpstreamPrivacyHeaders(headers)

	if _, ok := headers["x-real-ip"]; ok {
		t.Fatalf("expected x-real-ip to be removed")
	}
	if _, ok := headers["X-Original-URL"]; ok {
		t.Fatalf("expected X-Original-URL to be removed")
	}
	if _, ok := headers["X-Original-Method"]; ok {
		t.Fatalf("expected X-Original-Method to be removed")
	}
	if _, ok := headers["X-Forwarded-URI"]; ok {
		t.Fatalf("expected X-Forwarded-URI to be removed")
	}
	if _, ok := headers["X-Forwarded-Scheme"]; ok {
		t.Fatalf("expected X-Forwarded-Scheme to be removed")
	}
	if _, ok := headers["X-Forwarded-Method"]; ok {
		t.Fatalf("expected X-Forwarded-Method to be removed")
	}

	if got := headers.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type to be preserved, got %q", got)
	}
	if got := headers.Get("X-Custom"); got != "1" {
		t.Fatalf("expected X-Custom to be preserved, got %q", got)
	}
}

