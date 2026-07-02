package initialize

import (
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
)

// AppService for system contains all application services that are used in the system.
type AppService struct {
	// Domain repository

	// Application service

	// ...
}

// initializeAppService initializes the application services and returns an instance of AppService.
func initializeAppService(
	config *domain_config.SystemConfig,
	logger domain_logger.ILogger,
	connections *Connection,
) (*AppService, error) {
	// TODO: implement the initialization of application services and domain repositories here.
	logger.Info("Initializing all application services...")
	return &AppService{
		// Domain repository

		// Application service

		// ...
	}, nil
}
