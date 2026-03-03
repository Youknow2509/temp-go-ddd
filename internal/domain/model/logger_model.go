package model

// Logger config pkg model input
type LoggerConfigPkg struct {
	ServerMode string
	Setting    LoggerSetting
}

// Function newLoggerConfigPkg creates a new LoggerConfigPkg instance
func NewLoggerConfigPkg(serverMode string, setting LoggerSetting) *LoggerConfigPkg {
	return &LoggerConfigPkg{
		ServerMode: serverMode,
		Setting:    setting,
	}
}