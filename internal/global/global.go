package global

import (
	"sync"

	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
)

// Global variables for system
var (
	// Logger is the global variable to store logger instance
	Logger domain_logger.ILogger
	// SystemConfig is the global variable to store system configuration
	SystemConfig *domain_config.SystemConfig
	// WaitGroup is used to wait for all goroutines to finish before exiting the program
	WaitGroup *sync.WaitGroup
)
