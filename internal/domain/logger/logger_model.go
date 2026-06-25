package logger

import (
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// Logger config pkg model input
type LoggerConfigPkg struct {
	Setting domain_config.TelemetryLoggerSetting
}

// Function newLoggerConfigPkg creates a new LoggerConfigPkg instance
func NewLoggerConfigPkg(setting domain_config.TelemetryLoggerSetting) *LoggerConfigPkg {
	return &LoggerConfigPkg{
		Setting: setting,
	}
}
