package utils

import (
	"net/http"
	"strings"
)

// StripForwardedIPHeaders removes headers that may leak the original client IP
// (or proxy chain info) to upstream services.
//
// This is intentionally conservative: gpt-load acts as a proxy/load-balancer and
// should not forward client network identity headers by default.
func StripForwardedIPHeaders(header http.Header) {
	if header == nil {
		return
	}

	for key := range header {
		lowerKey := strings.ToLower(key)

		// Drop the whole X-Forwarded-* family (common reverse-proxy headers).
		if lowerKey == "x-forwarded" ||
			strings.HasPrefix(lowerKey, "x-forwarded-") ||
			strings.HasPrefix(lowerKey, "x_forwarded_") {
			header.Del(key)
			continue
		}

		// Standard / widely-used client IP headers.
		switch lowerKey {
		case "forwarded", // RFC 7239
			"client-ip",
			"true-client-ip",
			"x-true-client-ip",
			"x-real-ip",
			"x-real-proto"
			"x-client-ip",
			"x-cluster-client-ip",
			"x-forwarded-for", // covered by prefix, kept for clarity
			"x-http-forwarded-for",
			"x-originating-ip",
			"x-remote-addr",
			"x-remote-ip",
			"x-original-forwarded-for",
			"x-original-client-ip",
			"x-original-remote-addr",
			"x-proxyuser-ip",
			"x-appengine-user-ip",
			"x-envoy-external-address",
			"x-source-ip",
			"x-source-addr",
			"cf-connecting-ip",
			"fastly-client-ip",
			"fly-client-ip":
			header.Del(key)
		}
	}
}
