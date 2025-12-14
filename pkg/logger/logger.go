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

// Config 日志配置结构体
type Config struct {
	Level      string // 日志级别: debug, info, warn, error, fatal
	OutputPath string // 日志文件路径（为空则仅输出到控制台）
	MaxSize    int    // 日志文件最大大小（MB），超过后自动轮转
	MaxBackups int    // 保留的旧日志文件最大数量
	MaxAge     int    // 保留旧日志文件的最大天数
	Compress   bool   // 是否压缩轮转的日志文件
	Console    bool   // 是否输出到控制台
}

// DefaultConfig 返回默认的日志配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		OutputPath: "",
		MaxSize:    100,  // 100 MB
		MaxBackups: 3,    // 保留3个旧文件
		MaxAge:     28,   // 保留28天
		Compress:   true, // 压缩旧文件
		Console:    true, // 输出到控制台
	}
}

// Init 使用给定的配置初始化全局日志器
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// 创建编码器配置
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

	// 为不同的输出创建cores
	var cores []zapcore.Core

	// 使用lumberjack实现文件输出和轮转
	if config.OutputPath != "" {
		// 确保目录存在
		dir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.OutputPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		// 保存文件写入器以便后续关闭
		loggerMutex.Lock()
		fileWriterCloser = fileWriter
		loggerMutex.Unlock()

		// 文件输出使用JSON编码器（无颜色代码）
		fileEncoderConfig := encoderConfig
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出
	if config.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 如果没有配置输出，默认输出到控制台
	if len(cores) == 0 {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 使用多个cores创建日志器
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	loggerMutex.Lock()
	globalLogger = logger
	loggerMutex.Unlock()

	return nil
}

// InitSimple 使用简单参数初始化日志器（向后兼容）
func InitSimple(level string, outputPaths []string) error {
	config := DefaultConfig()
	config.Level = level
	if len(outputPaths) > 0 && outputPaths[0] != "" && outputPaths[0] != "stdout" {
		config.OutputPath = outputPaths[0]
	}
	return Init(config)
}

// Get 返回全局日志器
func Get() *zap.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if globalLogger == nil {
		// 如果未初始化，使用默认配置初始化
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// Info 记录info级别的日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error 记录error级别的日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn 记录warn级别的日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug 记录debug级别的日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal 记录fatal级别的日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync 刷新所有缓冲的日志条目
func Sync() error {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Close 关闭日志器并释放资源
func Close() error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	var err error
	if globalLogger != nil {
		err = globalLogger.Sync()
		globalLogger = nil
	}

	// 如果文件写入器存在，关闭它
	if fileWriterCloser != nil {
		if closeErr := fileWriterCloser.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		fileWriterCloser = nil
	}

	return err
}
