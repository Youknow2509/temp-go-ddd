package global

import (
	"sync"

	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	domain_model "github.com/youknow2509/temp-go-ddd/internal/domain/model"
)

// Global variables for system
var (
	// Logger is the global variable to store logger instance
	Logger domain_logger.ILogger
	// SystemConfig is the global variable to store system configuration
	SystemConfig *domain_model.SystemConfig
	// WaitGroup is used to wait for all goroutines to finish before exiting the program
	WaitGroup *sync.WaitGroup
)
