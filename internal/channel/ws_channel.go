package channel

import (
	"context"
	"net/http"

	"gpt-load/internal/models"

	"github.com/gin-gonic/gin"
)

// WSChannel is a generic WebSocket channel type. It reuses BaseChannel routing
// but does not modify requests or perform key validation.
type WSChannel struct {
	*BaseChannel
}

func init() {
	Register("ws", newWSChannel)
}

func newWSChannel(f *Factory, group *models.Group) (ChannelProxy, error) {
	base, err := f.newBaseChannel("ws", group)
	if err != nil {
		return nil, err
	}

	return &WSChannel{BaseChannel: base}, nil
}

// ModifyRequest is a no-op for WS channels; headers are forwarded as-is.
func (ch *WSChannel) ModifyRequest(req *http.Request, apiKey *models.APIKey, group *models.Group) {
}

// IsStreamRequest always returns false for WS channels in HTTP mode.
func (ch *WSChannel) IsStreamRequest(c *gin.Context, bodyBytes []byte) bool {
	return false
}

// ExtractModel is not applicable for WS channels.
func (ch *WSChannel) ExtractModel(c *gin.Context, bodyBytes []byte) string {
	return ""
}

// ValidateKey is a no-op for WS channels and always reports success.
func (ch *WSChannel) ValidateKey(ctx context.Context, apiKey *models.APIKey, group *models.Group) (bool, error) {
	return true, nil
}
