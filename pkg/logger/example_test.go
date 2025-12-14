package logger_test

import (
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"go.uber.org/zap"
)

// Example demonstrates basic logger usage with console output
func Example_basic() {
	// Initialize logger with default config (console output, info level)
	config := logger.DefaultConfig()
	logger.Init(config)
	defer logger.Close()

	// Log messages at different levels
	logger.Debug("This is a debug message")
	logger.Info("Application started", zap.String("version", "1.0.0"))
	logger.Warn("This is a warning")
	logger.Error("An error occurred", zap.Error(nil))
}

// Example demonstrates logger with file output and rotation
func Example_fileOutput() {
	// Configure logger with file output
	config := &logger.Config{
		Level:      "info",
		OutputPath: "/var/log/seatunnel/app.log",
		MaxSize:    100,  // 100 MB
		MaxBackups: 3,    // Keep 3 old files
		MaxAge:     28,   // Keep for 28 days
		Compress:   true, // Compress old files
		Console:    false,
	}

	logger.Init(config)
	defer logger.Close()

	logger.Info("Logging to file with rotation enabled")
}

// Example demonstrates logger with both console and file output
func Example_bothOutputs() {
	// Configure logger with both console and file output
	config := &logger.Config{
		Level:      "debug",
		OutputPath: "/var/log/seatunnel/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		Console:    true, // Also output to console
	}

	logger.Init(config)
	defer logger.Close()

	logger.Debug("This appears in both console and file")
	logger.Info("Application event", zap.String("event", "startup"))
}

// Example demonstrates simple initialization (backward compatibility)
func Example_simple() {
	// Simple initialization with just level and output path
	logger.InitSimple("info", []string{"/var/log/seatunnel/app.log"})
	defer logger.Close()

	logger.Info("Simple logger initialization")
}
