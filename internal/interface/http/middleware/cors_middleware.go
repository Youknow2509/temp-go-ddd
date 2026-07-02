package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// Cors Middleware
type CorsMiddleware struct {
	config domain_config.HttpCorsSetting
}

// NewCorsMiddleware creates a new instance of CorsMiddleware with the provided configuration.
func NewCorsMiddleware(config domain_config.HttpCorsSetting) *CorsMiddleware {
	return &CorsMiddleware{
		config: config,
	}
}

// Handle is the middleware function that sets the CORS headers for incoming HTTP requests.
func (m *CorsMiddleware) Handle() gin.HandlerFunc {
	cfg := m.config
	// If cors is disabled, allow all origins, methods, and headers
	if !cfg.Enabled {
		return cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}

	// Build the CORS configuration based on the provided settings
	corsConfig := cors.Config{
		AllowCredentials:    cfg.AllowCredentials,
		MaxAge:              time.Duration(cfg.MaxAgeSecs) * time.Second,
		AllowPrivateNetwork: cfg.AllowPrivateNetwork,
	}

	if cfg.Origin.AllowAnyOrigin {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = cfg.Origin.AllowedOrigins
	}

	if cfg.Methods.AllowAnyMethod {
		corsConfig.AllowMethods = []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		}
	} else {
		corsConfig.AllowMethods = cfg.Methods.AllowedMethods
	}

	if cfg.Headers.AllowAnyHeader {
		corsConfig.AllowHeaders = []string{"*"}
	} else {
		corsConfig.AllowHeaders = cfg.Headers.AllowedHeaders
	}

	corsConfig.ExposeHeaders = cfg.Headers.ExposedHeaders

	return cors.New(corsConfig)
}
