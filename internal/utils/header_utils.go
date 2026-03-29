package utils

import (
	"gpt-load/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var upstreamPrivacyHeaderKeys = []string{
	"Forwarded",
	"X-Forwarded-For",
	"X-Forwarded",
	"X-Forwarded-Host",
	"X-Forwarded-Port",
	"X-Forwarded-Proto",
	"X-Forwarded-Protocol",
	"X-Forwarded-Scheme",
	"X-Forwarded-Server",
	"X-Forwarded-Uri",
	"X-Forwarded-Method",
	"X-Original-Forwarded-For",
	"X-Original-Host",
	"X-Original-Method",
	"X-Original-Proto",
	"X-Original-Scheme",
	"X-Original-Uri",
	"X-Original-Url",
	"X-Real-IP",
	"X-Client-IP",
	"X-Cluster-Client-IP",
	"X-Originating-IP",
	"X-Remote-IP",
	"X-Remote-Addr",
	"True-Client-IP",
	"CF-Connecting-IP",
	"CF-Visitor",
	"CF-Ray",
	"CF-Country",
	"CDN-Loop",
	// Nginx control header; should not be forwarded to upstreams.
	"X-Accel-Buffering",
}

var upstreamPrivacyHeaderKeySet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(upstreamPrivacyHeaderKeys))
	for _, key := range upstreamPrivacyHeaderKeys {
		set[strings.ToLower(key)] = struct{}{}
	}
	return set
}()

// HeaderVariableContext holds context data for variable resolution
type HeaderVariableContext struct {
	ClientIP string
	Group    *models.Group
	APIKey   *models.APIKey
}

// ResolveHeaderVariables resolves dynamic variables in header values
func ResolveHeaderVariables(value string, ctx *HeaderVariableContext) string {
	if ctx == nil {
		return value
	}

	now := time.Now()
	result := value

	// Replace all supported variables
	variables := map[string]string{
		"${CLIENT_IP}":    ctx.ClientIP,
		"${TIMESTAMP_MS}": strconv.FormatInt(now.UnixMilli(), 10),
		"${TIMESTAMP_S}":  strconv.FormatInt(now.Unix(), 10),
	}

	if ctx.Group != nil {
		variables["${GROUP_NAME}"] = ctx.Group.Name
	}

	if ctx.APIKey != nil {
		variables["${API_KEY}"] = ctx.APIKey.KeyValue
	}

	// Replace variables in the value
	for variable, replacement := range variables {
		result = strings.ReplaceAll(result, variable, replacement)
	}

	return result
}

// ApplyHeaderRules applies header rules to the HTTP request
func ApplyHeaderRules(req *http.Request, rules []models.HeaderRule, ctx *HeaderVariableContext) {
	if req == nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		canonicalKey := http.CanonicalHeaderKey(rule.Key)

		switch rule.Action {
		case "remove":
			req.Header.Del(canonicalKey)
		case "set":
			resolvedValue := ResolveHeaderVariables(rule.Value, ctx)
			req.Header.Set(canonicalKey, resolvedValue)
		}
	}
}

// NewHeaderVariableContextFromGin creates HeaderVariableContext from Gin context
func NewHeaderVariableContextFromGin(c *gin.Context, group *models.Group, apiKey *models.APIKey) *HeaderVariableContext {
	if c == nil {
		return nil
	}

	return &HeaderVariableContext{
		ClientIP: c.ClientIP(),
		Group:    group,
		APIKey:   apiKey,
	}
}

// NewHeaderVariableContext creates HeaderVariableContext without Gin context
func NewHeaderVariableContext(group *models.Group, apiKey *models.APIKey) *HeaderVariableContext {
	return &HeaderVariableContext{
		ClientIP: "127.0.0.1",
		Group:    group,
		APIKey:   apiKey,
	}
}

// RemoveClientIPHeaders deletes headers that may leak the original client IP to upstream.
//
// Deprecated: use RemoveUpstreamPrivacyHeaders.
func RemoveClientIPHeaders(headers http.Header) {
	RemoveUpstreamPrivacyHeaders(headers)
}

// RemoveUpstreamPrivacyHeaders deletes headers that may leak client IP or internal routing info to upstream.
func RemoveUpstreamPrivacyHeaders(headers http.Header) {
	if headers == nil {
		return
	}

	for key := range headers {
		if _, ok := upstreamPrivacyHeaderKeySet[strings.ToLower(key)]; ok {
			delete(headers, key)
		}
	}
}
