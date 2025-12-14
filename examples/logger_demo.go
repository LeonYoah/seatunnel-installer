package main

import (
	"time"

	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Example 1: Console-only logging
	println("=== Example 1: Console-only logging ===")
	config1 := &logger.Config{
		Level:   "debug",
		Console: true,
	}
	logger.Init(config1)

	logger.Debug("Debug message - detailed information")
	logger.Info("Info message - general information", zap.String("app", "seatunnel"))
	logger.Warn("Warning message - something to watch")
	logger.Error("Error message - something went wrong", zap.Int("code", 500))

	logger.Close()
	time.Sleep(100 * time.Millisecond)

	// Example 2: File logging with rotation
	println("\n=== Example 2: File logging with rotation ===")
	config2 := &logger.Config{
		Level:      "info",
		OutputPath: "logs/app.log",
		MaxSize:    10,   // 10 MB
		MaxBackups: 3,    // Keep 3 old files
		MaxAge:     7,    // Keep for 7 days
		Compress:   true, // Compress old files
		Console:    false,
	}
	logger.Init(config2)

	logger.Info("Application started", zap.String("version", "1.0.0"))
	logger.Info("Configuration loaded", zap.Int("workers", 4))
	logger.Warn("High memory usage detected", zap.Float64("usage_percent", 85.5))

	logger.Close()
	println("Logs written to logs/app.log")
	time.Sleep(100 * time.Millisecond)

	// Example 3: Both console and file output
	println("\n=== Example 3: Both console and file output ===")
	config3 := &logger.Config{
		Level:      "debug",
		OutputPath: "logs/combined.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     7,
		Compress:   false,
		Console:    true,
	}
	logger.Init(config3)

	logger.Debug("This appears in both console and file")
	logger.Info("Processing task",
		zap.String("task_id", "task-123"),
		zap.String("status", "running"),
		zap.Duration("elapsed", 5*time.Second),
	)
	logger.Info("Task completed successfully",
		zap.String("task_id", "task-123"),
		zap.Int("records_processed", 1000),
	)

	logger.Close()
	println("\nDemo completed!")
}
