package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/youknow2509/temp-go-ddd/internal/constant"
	domain_config "github.com/youknow2509/temp-go-ddd/internal/domain/config"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ===
// Use go.uber.org/zap deployment logger
// Example usage:
// - logger.Info("This is an info message", "user_id", 123, "action", "login")
// - logger.Info(zap.String("user_id", "123"), zap.String("action", "login"))
// Structure log key-value pairs or zap.Fields with helperConvertToZapFields for better type handling.
// ===
type ZapLogger struct {
	logger *zap.Logger
}

// Debug implements [logger.ILogger].
func (z ZapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Debug(msg, helperConvertToZapFields(fields...)...)
}

// Error implements [logger.ILogger].
func (z ZapLogger) Error(msg string, fields ...interface{}) {
	z.logger.Error(msg, helperConvertToZapFields(fields...)...)
}

// Fatal implements [logger.ILogger].
func (z ZapLogger) Fatal(msg string, fields ...interface{}) {
	z.logger.Fatal(msg, helperConvertToZapFields(fields...)...)
}

// Info implements [logger.ILogger].
func (z ZapLogger) Info(msg string, fields ...interface{}) {
	z.logger.Info(msg, helperConvertToZapFields(fields...)...)
}

// Panic implements [logger.ILogger].
func (z ZapLogger) Panic(msg string, fields ...interface{}) {
	z.logger.Panic(msg, helperConvertToZapFields(fields...)...)
}

// Warn implements [logger.ILogger].
func (z ZapLogger) Warn(msg string, fields ...interface{}) {
	z.logger.Warn(msg, helperConvertToZapFields(fields...)...)
}

// New intance and implementation of logger interface
func NewZapLogger(config *domain_logger.LoggerConfigPkg) (domain_logger.ILogger, error) {
	if config == nil {
		return nil, fmt.Errorf("logger config is nil")
	}

	// Parse base log level
	baseLevel := helperParseZapLevel(config.Setting.Level)

	// Define encoders
	jsonEncoder := helperGetZapEncoderJson()
	consoleEncoder := helperGetZapEncoderConsole()

	var cores []zapcore.Core

	// Create outputs dynamically based on config.Setting.Output
	for _, out := range config.Setting.Output {
		switch out {
		case "stdout":
			cores = append(cores, zapcore.NewCore(
				consoleEncoder,
				zapcore.Lock(os.Stdout),
				baseLevel,
			))
		case "stderr":
			cores = append(cores, zapcore.NewCore(
				consoleEncoder,
				zapcore.Lock(os.Stderr),
				baseLevel,
			))
		case "file":
			if config.Setting.File.Enabled {
				fileWriter := helperGetFileWriter(config.Setting)
				cores = append(cores, zapcore.NewCore(
					jsonEncoder,
					fileWriter,
					baseLevel,
				))
			}
		}
	}

	// Fallback to stdout if no core is created
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			baseLevel,
		))
	}

	// Tee all cores together
	core := zapcore.NewTee(cores...)

	// Build zap options dynamically
	var zapOptions []zap.Option

	if config.Setting.Caller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if config.Setting.StacktraceLevel != "" {
		stacktraceLevel := helperParseZapLevel(config.Setting.StacktraceLevel)
		zapOptions = append(zapOptions, zap.AddStacktrace(stacktraceLevel))
	}

	// Create Zap Logger
	logger := zap.New(core, zapOptions...)

	return ZapLogger{
		logger: logger,
	}, nil
}

// ===
// Helper function for zap logger
// ===

// Helper parse string log level to zapcore.Level
func helperParseZapLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "trace", "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Helper convert variadic interface{} to zap.Field slice
//   - Accepts a variadic list of fields, which can be either zap.Field or key-value pairs (string followed by any value).
//   - Examples of usage:
//     helperConvertToZapFields("user_id", 123, "action", "login")
//     helperConvertToZapFields(zap.String("user_id", "123"), zap.String("action", "login"))
func helperConvertToZapFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)/2)

	i := 0
	for i < len(fields) {
		// Nếu field đã là zap.Field -> giữ nguyên
		if f, ok := fields[i].(zap.Field); ok {
			zapFields = append(zapFields, f)
			i++
			continue
		}

		// Lấy key (phải là string)
		key, ok := fields[i].(string)
		if !ok {
			zapFields = append(zapFields, zap.Any("unknown", fields[i]))
			i++
			continue
		}

		// Kiểm tra có value đi kèm không
		if i+1 >= len(fields) {
			zapFields = append(zapFields, zap.String(key, ""))
			i++
			continue
		}

		// Lấy value và chuyển đổi theo kiểu
		value := fields[i+1]
		switch v := value.(type) {
		// Primitive types - string, int, float, bool...
		case string:
			zapFields = append(zapFields, zap.String(key, v))
		case int:
			zapFields = append(zapFields, zap.Int(key, v))
		case int8:
			zapFields = append(zapFields, zap.Int8(key, v))
		case int16:
			zapFields = append(zapFields, zap.Int16(key, v))
		case int32:
			zapFields = append(zapFields, zap.Int32(key, v))
		case int64:
			zapFields = append(zapFields, zap.Int64(key, v))
		case uint:
			zapFields = append(zapFields, zap.Uint(key, v))
		case uint8:
			zapFields = append(zapFields, zap.Uint8(key, v))
		case uint16:
			zapFields = append(zapFields, zap.Uint16(key, v))
		case uint32:
			zapFields = append(zapFields, zap.Uint32(key, v))
		case uint64:
			zapFields = append(zapFields, zap.Uint64(key, v))
		case float32:
			zapFields = append(zapFields, zap.Float32(key, v))
		case float64:
			zapFields = append(zapFields, zap.Float64(key, v))
		case bool:
			zapFields = append(zapFields, zap.Bool(key, v))
		// Time types
		case time.Time:
			zapFields = append(zapFields, zap.Time(key, v))
		case time.Duration:
			zapFields = append(zapFields, zap.Duration(key, v))
		// Error type
		case error:
			zapFields = append(zapFields, zap.NamedError(key, v))
		// Stringer interface (fmt.Stringer)
		// Bất kỳ object nào implement String() string
		case fmt.Stringer:
			zapFields = append(zapFields, zap.Stringer(key, v))
		// Object types: struct, map, slice, pointer...
		// zap.Any sử dụng JSON marshaling -> hỗ trợ tất cả
		default:
			zapFields = append(zapFields, zap.Any(key, v))
		}
		i += 2
	}

	return zapFields
}

// Helper get file writer with lumberjack rotation
func helperGetFileWriter(setting domain_config.TelemetryLoggerSetting) zapcore.WriteSyncer {
	// Create log directory if not exists
	if _, err := os.Stat(setting.File.Folder); os.IsNotExist(err) {
		os.MkdirAll(setting.File.Folder, os.ModePerm)
	}
	// Create lumberjack logger for file rotation
	lumberJackLogger := &lumberjack.Logger{
		Filename:   setting.File.Folder + "/" + setting.File.Filename, // Đường dẫn file log
		MaxSize:    setting.File.MaxSizeMb,                            // MB trước khi rotate
		MaxBackups: setting.File.MaxBackups,                           // Số file backup giữ lại
		MaxAge:     setting.File.MaxAgeDays,                           // Số ngày giữ file cũ
		Compress:   setting.File.Compress,                             // Nén file log cũ (gzip)
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Helper get zap encoder console
func helperGetZapEncoderConsole() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// Helper get zap encoder json
func helperGetZapEncoderJson() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	// Override the default keys with our custom keys from constant package
	cfg.LevelKey = constant.LoggerKeyLevelMsg
	cfg.TimeKey = constant.LoggerKeyCurrentTime
	cfg.CallerKey = constant.LoggerKeyCaller
	cfg.MessageKey = constant.LoggerKeyMessage
	cfg.StacktraceKey = constant.LoggerKeyStacktrace
	// Return a new JSON encoder with the custom configuration
	return zapcore.NewJSONEncoder(cfg)
}
