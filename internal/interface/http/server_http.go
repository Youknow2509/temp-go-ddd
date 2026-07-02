package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/youknow2509/temp-go-ddd/internal/constant"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	middleware "github.com/youknow2509/temp-go-ddd/internal/interface/http/middleware"
)

// Server bundles a gin.Engine with an *http.Server configured entirely from config
type Server struct {
	cfg        *domain_config.HttpServerSetting
	engine     *gin.Engine
	httpServer *http.Server
	logger     domain_logger.ILogger
	waitGroup  *sync.WaitGroup
}

// New creates a new HTTP Server instance configured from configuration settings.
func New(
	config domain_config.HttpServerSetting,
	logger domain_logger.ILogger,
	registerRoutes func(*gin.Engine),
) (*Server, error) {
	// Get mode server from env
	server_mode := os.Getenv(constant.SystemModeEnvKey)
	var engine *gin.Engine
	switch server_mode {
	case constant.SystemModeDevelopment:
		engine = gin.Default()
		engine.Use(gin.Recovery())
		engine.Use(gin.Logger())
		engine.Use(gin.ErrorLogger())
	default:
		// Handle production mode
		engine = gin.New()
	}
	// Setup cors middleware
	cors_middleware := middleware.NewCorsMiddleware(config.Security.Cors)
	engine.Use(cors_middleware.Handle())
	// Setup rate limit middleware
	rate_limit_middleware := middleware.NewRateLimitMiddleware(config.RateLimit)
	engine.Use(rate_limit_middleware.Handle())
	// Setup body size limit middleware
	body_size_limit_middleware := middleware.NewBodyLimitMiddleware(config.Limits)
	engine.Use(body_size_limit_middleware.Handle())

	// Register routes
	if registerRoutes != nil {
		registerRoutes(engine)
	}

	// Build addr for the server
	addr := fmt.Sprintf("%s:%d", config.Network.Host, config.Network.Port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
		// timeouts
		ReadTimeout:  time.Duration(config.Timeouts.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(config.Timeouts.WriteTimeoutMs) * time.Millisecond,
		IdleTimeout:  time.Duration(config.Timeouts.IdleTimeoutMs) * time.Millisecond,
		// limits.max_header_size
		MaxHeaderBytes: config.Limits.MaxHeaderSize,
	}

	// security.tls
	if config.Security.Tls.IsEnabled {
		tlsCfg, err := buildTLSConfig(config.Security.Tls)
		if err != nil {
			logger.Error("build tls config: ", err)
			return nil, fmt.Errorf("build tls config: %w", err)
		}
		httpServer.TLSConfig = tlsCfg
	}

	return &Server{
		cfg:        &config,
		engine:     engine,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

// Start binds the listener (applying worker/tcp settings) and serves
// until Shutdown is called. It blocks the calling goroutine.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.logger.Error("listen on ", s.httpServer.Addr, ": ", err)
		return fmt.Errorf("listen on %s: %w", s.httpServer.Addr, err)
	}

	tcpLn, ok := ln.(*net.TCPListener)
	if !ok {
		s.logger.Error("listen on ", s.httpServer.Addr, ": expected *net.TCPListener, got %T", ln)
		return fmt.Errorf("expected *net.TCPListener, got %T", ln)
	}

	// tcp.*
	tuned := newTCPTunedListener(tcpLn, s.cfg.Tcp)

	s.logger.Info(
		"http server listening on ", s.httpServer.Addr,
		" tls=", s.cfg.Security.Tls.IsEnabled,
	)

	if s.cfg.Security.Tls.IsEnabled {
		err = s.httpServer.ServeTLS(tuned, s.cfg.Security.Tls.CertFile, s.cfg.Security.Tls.KeyFile)
	} else {
		err = s.httpServer.Serve(tuned)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("http server error: ", err)
		return err
	}
	return nil
}

// Shutdown gracefully drains in-flight requests, bounded by timeouts.shutdown_timeout_ms.
func (s *Server) Shutdown() error {
	timeout := time.Duration(s.cfg.Timeouts.ShutdownTimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.logger.Info("shutting down http server...")
	return s.httpServer.Shutdown(ctx)
}
