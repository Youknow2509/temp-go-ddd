package initialize

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
)

// WaitForShutdown listens for OS signals to gracefully shut down the application
func WaitForShutdown(cancel context.CancelFunc, wg *sync.WaitGroup, res *AppResources) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-quit

	// Determine logger to use
	var logger domain_logger.ILogger
	if res != nil {
		logger = res.Logger
	}

	logMsg := func(msg string, err error) {
		if logger != nil {
			if err != nil {
				logger.Error(msg, "error", err)
			} else {
				logger.Info(msg)
			}
		} else {
			if err != nil {
				log.Printf("[ERROR] %s: %v\n", msg, err)
			} else {
				log.Printf("[INFO] %s\n", msg)
			}
		}
	}

	logMsg(fmt.Sprintf("System shutdown signal received: %s", sig.String()), nil)

	// Cancel root context to notify background routines to stop
	if cancel != nil {
		cancel()
	}

	// Create a context with a timeout for shutdown operations
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if res != nil {
		// Gracefully shut down Metrics HTTP Server
		if res.MetricServer != nil {
			logMsg("Shutting down metrics HTTP server ...", nil)
			if err := res.MetricServer.Shutdown(ctx); err != nil {
				logMsg("Metrics HTTP server shutdown error", err)
			} else {
				logMsg("Metrics HTTP server shut down successfully", nil)
			}
		}

		// Gracefully shut down TracerProvider (flush remaining traces)
		if res.TracerProvider != nil {
			logMsg("Shutting down tracer provider ...", nil)
			if err := res.TracerProvider.Shutdown(ctx); err != nil {
				logMsg("Tracer provider shutdown error", err)
			} else {
				logMsg("Tracer provider shut down successfully", nil)
			}
		}

		// Gracefully shut down MeterProvider
		if res.MeterProvider != nil {
			logMsg("Shutting down meter provider ...", nil)
			if err := res.MeterProvider.Shutdown(ctx); err != nil {
				logMsg("Meter provider shutdown error", err)
			} else {
				logMsg("Meter provider shut down successfully", nil)
			}
		}

		// Stop all interfaces http, grpc, ws, kafka, ...
		shutdownInterfaces(res, logMsg)

		// Stop all connections to external services (e.g., databases, message brokers, etc.)
		shutdownConnections(res, logMsg)

	}

	// Wait for other goroutines in WaitGroup to finish
	if wg != nil {
		logMsg("Waiting for background goroutines to finish ...", nil)
		wg.Wait()
	}

	logMsg("Graceful shutdown completed. Exiting system.", nil)
}

// shutdownInterfaces gracefully shuts down all interfaces (e.g., HTTP, gRPC, WebSocket, Kafka, etc.)
func shutdownInterfaces(res *AppResources, logMsg func(string, error)) {
	if res.Interfaces == nil {
		return
	}
	logMsg("Shutting down all interfaces", nil)

	// Shutdown HTTP server
	if res.Interfaces.HttpServer != nil {
		logMsg("Shutting down HTTP server", nil)
		if err := res.Interfaces.HttpServer.Shutdown(); err != nil {
			logMsg("HTTP server shutdown error", err)
		} else {
			logMsg("HTTP server shut down successfully", nil)
		}
	}

	// Shutdown gRPC server
	if res.Interfaces.GrpcServer != nil {
		logMsg("Shutting down gRPC server", nil)
		res.Interfaces.GrpcServer.GracefulStop()
		logMsg("gRPC server shut down successfully", nil)
	}
}

// shutdownConnections gracefully closes all external connection pools/clients
func shutdownConnections(res *AppResources, logMsg func(string, error)) {
	if res.Connections == nil {
		return
	}
	logMsg("Shutting down connections to external services", nil)

	if res.Connections.Postgres != nil {
		logMsg("Closing Postgres connection pool", nil)
		res.Connections.Postgres.Close()
	}

	if res.Connections.Redis != nil {
		logMsg("Closing Redis client", nil)
		if err := res.Connections.Redis.Close(); err != nil {
			logMsg("Redis client close error", err)
		} else {
			logMsg("Redis client closed successfully", nil)
		}
	}

	if res.Connections.ScyllaDb != nil {
		logMsg("Closing ScyllaDB session", nil)
		res.Connections.ScyllaDb.Close()
	}

	if res.Connections.Kafka != nil {
		logMsg("Closing Kafka producers", nil)
		if err := res.Connections.Kafka.Close(); err != nil {
			logMsg("Kafka producers close error", err)
		} else {
			logMsg("Kafka producers closed successfully", nil)
		}
	}

	for name, client := range res.Connections.GrpcClients {
		if client != nil {
			logMsg(fmt.Sprintf("Closing gRPC client connection for '%s'", name), nil)
			if err := client.Close(); err != nil {
				logMsg(fmt.Sprintf("gRPC client connection for '%s' close error", name), err)
			} else {
				logMsg(fmt.Sprintf("gRPC client connection for '%s' closed successfully", name), nil)
			}
		}
	}
}


