package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/youknow2509/temp-go-ddd/internal/constant"
	domain_logger "github.com/youknow2509/temp-go-ddd/internal/domain/logger"
	domain_model "github.com/youknow2509/temp-go-ddd/internal/domain/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ==============================================================
// Use go.uber.org/zap deployment logger
// Example usage:
// - logger.Info("This is an info message", "user_id", 123, "action", "login")
// - logger.Info(zap.String("user_id", "123"), zap.String("action", "login"))
// Structure log key-value pairs or zap.Fields with helperConvertToZapFields for better type handling.
// ==============================================================
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
func NewZapLogger(config *domain_model.LoggerConfigPkg) (domain_logger.ILogger, error) {
	if config == nil {
		return nil, fmt.Errorf("logger config is nil")
	}
	// Define level handling logic
	logger_level := helperGetZapLevelEnabler(config.ServerMode)

	// Define encoders
	json_encoder := helperGetZapEncoderJson()
	console_encoder := helperGetZapEncoderConsole()

	// Define file writer with lumberjack rotation
	file_writer := helperGetFileWriter(config.Setting)

	// Build cores based on server mode
	var cores []zapcore.Core

	// Always log to file
	cores = append(cores, zapcore.NewCore(json_encoder, file_writer, logger_level))
	// Check mode to decide if we also log to console
	if config.ServerMode == constant.SystemModeDevelopment {
		console_writer := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(
			console_encoder,
			console_writer,
			logger_level,
		))
	}

	// Join cores together
	core := zapcore.NewTee(cores...)

	// Build zap options based on server mode
	zapOptions := helperGetZapOptions(config.ServerMode)

	// Create construct Logger.
	logger := zap.New(core, zapOptions...)

	return ZapLogger{
		logger: logger,
	}, nil
}

// ==============================================================
// Helper function for zap logger
// ==============================================================

// Helper get zap options by server mode
func helperGetZapOptions(server_mode string) []zap.Option {
	switch server_mode {
	case constant.SystemModeDevelopment:
		return []zap.Option{
			zap.AddCaller(),                      // Ghi file:line gọi log
			zap.AddCallerSkip(0),                 // Caller skip level (điều chỉnh nếu wrap thêm layer)
			zap.AddStacktrace(zapcore.WarnLevel), // Stacktrace từ Warn trở lên -> debug dễ
			zap.Development(),                    // DPanic sẽ panic thay vì chỉ log
		}
	default: // Production
		return []zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel), // Chỉ stacktrace khi Panic/Fatal -> giảm noise
		}
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
func helperGetFileWriter(setting domain_model.LoggerSetting) zapcore.WriteSyncer {
	// Create log directory if not exists
	if _, err := os.Stat(setting.FolderStore); os.IsNotExist(err) {
		os.MkdirAll(setting.FolderStore, os.ModePerm)
	}
	// Create lumberjack logger for file rotation
	lumberJackLogger := &lumberjack.Logger{
		Filename:   setting.FolderStore + "/app.log", // Đường dẫn file log
		MaxSize:    setting.FileMaxSize,              // MB trước khi rotate
		MaxBackups: setting.FileMaxBackups,           // Số file backup giữ lại
		MaxAge:     setting.FileMaxAge,               // Số ngày giữ file cũ
		Compress:   setting.Compress,                 // Nén file log cũ (gzip)
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

// Helper get zap level handler by server mode
func helperGetZapLevelEnabler(server_mode string) zap.LevelEnablerFunc {
	switch server_mode {
	case constant.SystemModeDevelopment:
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.WarnLevel
		})
	default:
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel
		})
	}
}
