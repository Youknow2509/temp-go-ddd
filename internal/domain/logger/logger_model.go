package logger

import (
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
)

// Logger config pkg model input
type LoggerConfigPkg struct {
	ServerMode string
	Setting    domain_config.LoggerSetting
}

// Function newLoggerConfigPkg creates a new LoggerConfigPkg instance
func NewLoggerConfigPkg(serverMode string, setting domain_config.LoggerSetting) *LoggerConfigPkg {
	return &LoggerConfigPkg{
		ServerMode: serverMode,
		Setting:    setting,
	}
}
