package logger

/**
 * Interface logger service
 */
type ILogger interface {
	// Info logs
	Info(msg string, fields ...interface{})
	// Error logs
	Error(msg string, fields ...interface{})
	// Warn logs
	Warn(msg string, fields ...interface{})
	// Panic logs
	Panic(msg string, fields ...interface{})
	// Debug logs
	Debug(msg string, fields ...interface{})
	// Fatal logs
	Fatal(msg string, fields ...interface{})
}

