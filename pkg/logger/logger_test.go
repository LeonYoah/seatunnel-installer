package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", config.Level)
	}
	if config.MaxSize != 100 {
		t.Errorf("Expected default MaxSize 100, got %d", config.MaxSize)
	}
	if config.MaxBackups != 3 {
		t.Errorf("Expected default MaxBackups 3, got %d", config.MaxBackups)
	}
	if !config.Console {
		t.Error("Expected Console to be true by default")
	}
}

func TestInitWithConsoleOnly(t *testing.T) {
	config := &Config{
		Level:   "debug",
		Console: true,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	logger := Get()
	if logger == nil {
		t.Fatal("Logger should not be nil after initialization")
	}

	// Test logging at different levels
	Debug("Debug message", zap.String("key", "value"))
	Info("Info message", zap.Int("count", 42))
	Warn("Warning message")
	Error("Error message", zap.Error(os.ErrNotExist))
}

func TestInitWithFileOutput(t *testing.T) {
	// Create temp directory for logs
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	config := &Config{
		Level:      "info",
		OutputPath: logFile,
		MaxSize:    1, // 1 MB for testing
		MaxBackups: 2,
		Console:    false,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger with file output: %v", err)
	}
	defer Close()

	// Write some log messages
	Info("Test message 1")
	Info("Test message 2", zap.String("key", "value"))
	Warn("Warning message")

	// Sync to ensure logs are written
	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before reading file
	Close()

	// Verify log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFile)
	}

	// Read log file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Test message 1") {
		t.Error("Log file should contain 'Test message 1'")
	}
	if !strings.Contains(contentStr, "Test message 2") {
		t.Error("Log file should contain 'Test message 2'")
	}
}

func TestInitWithBothOutputs(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "both.log")

	config := &Config{
		Level:      "debug",
		OutputPath: logFile,
		MaxSize:    10,
		MaxBackups: 3,
		Console:    true,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger with both outputs: %v", err)
	}
	defer Close()

	Debug("Debug to both outputs")
	Info("Info to both outputs")

	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before checking file
	Close()

	// Verify file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFile)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		level    string
		expected string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"invalid", "info"}, // Should default to info
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			config := &Config{
				Level:   tt.level,
				Console: false, // Disable console to avoid noise
			}

			err := Init(config)
			if err != nil {
				t.Fatalf("Failed to initialize logger with level %s: %v", tt.level, err)
			}

			logger := Get()
			if logger == nil {
				t.Fatal("Logger should not be nil")
			}
		})
	}
}

func TestLogRotation(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "rotate.log")

	// Create a logger with very small max size for testing rotation
	config := &Config{
		Level:      "info",
		OutputPath: logFile,
		MaxSize:    1, // 1 MB
		MaxBackups: 2,
		MaxAge:     7,
		Compress:   false, // Don't compress for easier testing
		Console:    false,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	// Write enough data to potentially trigger rotation
	data := strings.Repeat("This is a test log line that will help trigger rotation. ", 100)
	for i := 0; i < 100; i++ {
		Info(data, zap.Int("iteration", i))
	}

	// Sync to ensure all writes are flushed
	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before checking file
	Close()

	// Check if main log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Main log file should exist")
	}
}

func TestInitSimple(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "simple.log")

	err := InitSimple("info", []string{logFile})
	if err != nil {
		t.Fatalf("Failed to initialize logger with InitSimple: %v", err)
	}
	defer Close()

	Info("Simple init test")
	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before checking file
	Close()

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFile)
	}
}

func TestHelperFunctions(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "helpers.log")

	config := &Config{
		Level:      "debug",
		OutputPath: logFile,
		Console:    false,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	// Test all helper functions
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before reading file
	Close()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	expectedMessages := []string{"Debug message", "Info message", "Warn message", "Error message"}
	for _, msg := range expectedMessages {
		if !strings.Contains(contentStr, msg) {
			t.Errorf("Log file should contain '%s'", msg)
		}
	}
}

func TestGetBeforeInit(t *testing.T) {
	// Reset global logger
	loggerMutex.Lock()
	globalLogger = nil
	loggerMutex.Unlock()

	// Get should return a default logger even if not initialized
	logger := Get()
	if logger == nil {
		t.Fatal("Get() should return a default logger when not initialized")
	}

	// Should be able to log without errors
	Info("Test message before init")
}

func TestConcurrentLogging(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "concurrent.log")

	config := &Config{
		Level:      "info",
		OutputPath: logFile,
		MaxSize:    10,
		MaxBackups: 3,
		Console:    false,
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	// Spawn multiple goroutines to log concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				Info("Concurrent log", zap.Int("goroutine", id), zap.Int("iteration", j))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	if err := Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Close logger before checking file
	Close()

	// Verify log file exists and has content
	info, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("Log file should exist: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Log file should not be empty after concurrent logging")
	}
}
