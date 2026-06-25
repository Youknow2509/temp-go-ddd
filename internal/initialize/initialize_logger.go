package initialize

import (
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	"github.com/youknow2509/temp-go-ddd/internal/global"
	infra_logger "github.com/youknow2509/temp-go-ddd/internal/infrastructure/logger"
)

// ===
// private function to initialize the logger
// ===

func initializeLogger() error {
	// Create logger configuration
	logger_config := domain_logger.NewLoggerConfigPkg(
		global.SystemConfig.System.Mode,
		global.SystemConfig.Logger,
	)
	// Create logger instance
	logger, err := infra_logger.NewZapLogger(logger_config)
	if err != nil {
		return err
	}
	logger.Debug("Logger initialized ...")
	logger.Debug("System Config: ", "data", global.SystemConfig)
	// Set logger instance to global variable
	global.Logger = logger
	return nil
}
