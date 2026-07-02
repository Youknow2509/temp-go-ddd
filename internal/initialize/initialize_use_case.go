package initialize

import (
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
)

// initializeUseCase initializes the use cases.
func initializeUseCase(
	config *domain_config.SystemConfig,
	logger domain_logger.ILogger,
	services *AppService,
) error {
	// TODO: Implement the initialization of use cases here.
	logger.Info("Initializing all use cases...")
	return nil
}
