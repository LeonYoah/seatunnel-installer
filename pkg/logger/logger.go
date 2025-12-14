package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger     *zap.Logger
	fileWriterCloser io.Closer
	loggerMutex      sync.RWMutex
)

// Config holds logger configuration
type Config struct {
	Level      string // Log level: debug, info, warn, error, fatal
	OutputPath string // File path for log output (empty for console only)
	MaxSize    int    // Maximum size in megabytes before rotation
	MaxBackups int    // Maximum number of old log files to retain
	MaxAge     int    // Maximum number of days to retain old log files
	Compress   bool   // Whether to compress rotated log files
	Console    bool   // Whether to output to console
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		OutputPath: "",
		MaxSize:    100,  // 100 MB
		MaxBackups: 3,    // Keep 3 old files
		MaxAge:     28,   // Keep for 28 days
		Compress:   true, // Compress old files
		Console:    true, // Output to console
	}
}

// Init initializes the global logger with the given configuration
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create cores for different outputs
	var cores []zapcore.Core

	// File output with rotation using lumberjack
	if config.OutputPath != "" {
		// Ensure directory exists
		dir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.OutputPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		// Store the file writer for later closing
		loggerMutex.Lock()
		fileWriterCloser = fileWriter
		loggerMutex.Unlock()

		// Use JSON encoder for file output (no color codes)
		fileEncoderConfig := encoderConfig
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// Console output
	if config.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// If no output is configured, default to console
	if len(cores) == 0 {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// Create logger with multiple cores
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	loggerMutex.Lock()
	globalLogger = logger
	loggerMutex.Unlock()

	return nil
}

// InitSimple initializes the logger with simple parameters (backward compatibility)
func InitSimple(level string, outputPaths []string) error {
	config := DefaultConfig()
	config.Level = level
	if len(outputPaths) > 0 && outputPaths[0] != "" && outputPaths[0] != "stdout" {
		config.OutputPath = outputPaths[0]
	}
	return Init(config)
}

// Get returns the global logger
func Get() *zap.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if globalLogger == nil {
		// Initialize with default config if not initialized
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Close closes the logger and releases resources
func Close() error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	var err error
	if globalLogger != nil {
		err = globalLogger.Sync()
		globalLogger = nil
	}

	// Close the file writer if it exists
	if fileWriterCloser != nil {
		if closeErr := fileWriterCloser.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		fileWriterCloser = nil
	}

	return err
}
