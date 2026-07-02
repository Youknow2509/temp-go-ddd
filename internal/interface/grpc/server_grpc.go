package grpc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server bundles a *grpc.Server configured from config
type Server struct {
	cfg        *domain_config.GrpcServerSetting
	grpcServer *grpc.Server
	logger     domain_logger.ILogger
}

// New creates a new gRPC Server instance with configured keepalives, limits, TLS, and protocol options
func New(
	config domain_config.GrpcServerSetting,
	logger domain_logger.ILogger,
	registerServices func(*grpc.Server),
) (*Server, error) {
	var opts []grpc.ServerOption

	// Keepalive parameters
	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle: time.Duration(config.Timeouts.IdleTimeoutMs) * time.Millisecond,
		Time:              time.Duration(config.Keepalive.Http2KeepaliveIntervalMs) * time.Millisecond,
		Timeout:           time.Duration(config.Keepalive.Http2KeepaliveTimeoutMs) * time.Millisecond,
	}
	opts = append(opts, grpc.KeepaliveParams(keepaliveParams))

	// Keepalive enforcement policy
	enforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             time.Duration(config.Keepalive.Http2KeepaliveIntervalMs) * time.Millisecond,
		PermitWithoutStream: config.Keepalive.KeepaliveWhileIdle,
	}
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(enforcementPolicy))

	// Msg limits
	if config.Limits.MaxDecodingMessageSize > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(config.Limits.MaxDecodingMessageSize))
	}
	if config.Limits.MaxEncodingMessageSize > 0 {
		opts = append(opts, grpc.MaxSendMsgSize(config.Limits.MaxEncodingMessageSize))
	}

	// Http2 settings
	if config.Http2.InitialStreamWindowSize > 0 {
		opts = append(opts, grpc.InitialWindowSize(int32(config.Http2.InitialStreamWindowSize)))
	}
	if config.Http2.InitialConnectionWindowSize > 0 {
		opts = append(opts, grpc.InitialConnWindowSize(int32(config.Http2.InitialConnectionWindowSize)))
	}
	if config.Http2.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(uint32(config.Http2.MaxConcurrentStreams)))
	}

	// TLS Setup
	if config.Security.Tls.IsEnabled {
		cert, err := tls.LoadX509KeyPair(config.Security.Tls.CertFile, config.Security.Tls.KeyFile)
		if err != nil {
			logger.Error("load grpc tls key pair: ", err)
			return nil, fmt.Errorf("load x509 key pair: %w", err)
		}

		minVersion, err := parseTLSVersion(config.Security.Tls.MinVersion)
		if err != nil {
			logger.Error("parse grpc tls version: ", err)
			return nil, fmt.Errorf("parse tls version: %w", err)
		}

		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   minVersion,
		}

		if config.Security.Tls.RequireClientCert {
			tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		}

		creds := credentials.NewTLS(tlsCfg)
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	// Register reflection if enabled
	if config.Protocol.ReflectionEnabled {
		reflection.Register(grpcServer)
	}

	// Register health check if enabled
	if config.Protocol.HealthCheckEnabled {
		grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	}

	// Register other services
	if registerServices != nil {
		registerServices(grpcServer)
	}

	return &Server{
		cfg:        &config,
		grpcServer: grpcServer,
		logger:     logger,
	}, nil
}

// Start binds and starts serving gRPC requests
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Network.Host, s.cfg.Network.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("grpc listen on ", addr, " error: ", err)
		return fmt.Errorf("grpc listen on %s: %w", addr, err)
	}

	tcpLn, ok := ln.(*net.TCPListener)
	if !ok {
		s.logger.Error("grpc listen on ", addr, ": expected *net.TCPListener, got %T", ln)
		return fmt.Errorf("expected *net.TCPListener, got %T", ln)
	}

	tuned := newTCPTunedListener(tcpLn, s.cfg.Tcp)

	s.logger.Info(
		"grpc server listening on ", addr,
		" tls=", s.cfg.Security.Tls.IsEnabled,
	)

	if err := s.grpcServer.Serve(tuned); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		s.logger.Error("grpc server error: ", err)
		return err
	}
	return nil
}

// GracefulStop stops the gRPC server gracefully
func (s *Server) GracefulStop() {
	s.logger.Info("shutting down grpc server...")
	s.grpcServer.GracefulStop()
}

func parseTLSVersion(v string) (uint16, error) {
	switch v {
	case "", "1.2", "TLS1.2", "TLSv1.2":
		return tls.VersionTLS12, nil
	case "1.3", "TLS1.3", "TLSv1.3":
		return tls.VersionTLS13, nil
	case "1.1", "TLS1.1", "TLSv1.1":
		return tls.VersionTLS11, nil
	case "1.0", "TLS1.0", "TLSv1.0":
		return tls.VersionTLS10, nil
	default:
		return 0, fmt.Errorf("unsupported tls min_version: %q", v)
	}
}
