package middleware

import (
	"github.com/gin-gonic/gin"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// RateLimitMiddleware is a middleware
type RateLimitMiddleware struct {
	config domain_config.HttpRateLimitSetting
}

// NewRateLimitMiddleware creates a new RateLimitMiddleware
func NewRateLimitMiddleware(config domain_config.HttpRateLimitSetting) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		config: config,
	}
}

// Handle is a middleware handler
func (m *RateLimitMiddleware) Handle() gin.HandlerFunc {
	// TODO: Implement rate limiting logic based on m.config
	return func(c *gin.Context) {
		// Pre-handler phase
		c.Next()

		// Post-handler phase
	}
}
