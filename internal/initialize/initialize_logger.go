package initialize

import (
	domain_model "github.com/youknow2509/temp-go-ddd/internal/domain/model"
	"github.com/youknow2509/temp-go-ddd/internal/global"
	pkg_logger "github.com/youknow2509/temp-go-ddd/pkg/logger"
)

// ==============================================================
// private function to initialize the logger
// ==============================================================

func initializeLogger() error {
	// Create logger configuration
	logger_config := domain_model.NewLoggerConfigPkg(
		global.SystemConfig.System.Mode,
		global.SystemConfig.Logger,
	)
	// Create logger instance
	logger, err := pkg_logger.NewZapLogger(logger_config)
	if err != nil {
		return err
	}
	logger.Debug("Logger initialized ...")
	logger.Debug("System Config: ", "data", global.SystemConfig)
	// Set logger instance to global variable
	global.Logger = logger
	return nil
}
