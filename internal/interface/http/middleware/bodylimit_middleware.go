package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// BodyLimitMiddleware is a middleware that limits the size of the request body.
type BodyLimitMiddleware struct {
	config domain_config.HttpLimitsSetting
}

// NewBodyLimitMiddleware creates a new BodyLimitMiddleware with the given configuration.
func NewBodyLimitMiddleware(config domain_config.HttpLimitsSetting) *BodyLimitMiddleware {
	return &BodyLimitMiddleware{
		config: config,
	}
}

// Handler returns a gin.HandlerFunc that limits the size of the request body.
func (m *BodyLimitMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.config.MaxBodySize > 0 {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, m.config.MaxBodySize)
		}
		c.Next()
	}
}
