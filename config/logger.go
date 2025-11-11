package config

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/techmaster-vietnam/goerrorkit/core"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogrusLogger implement core.Logger interface sử dụng logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// Error implements core.Logger
func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}

// Info implements core.Logger
func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

// Debug implements core.Logger
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Debug(msg)
}

// Warn implements core.Logger
func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}

// LoggerOptions cấu hình cho logger
type LoggerOptions struct {
	// ConsoleOutput - Log ra console hay không
	ConsoleOutput bool

	// FileOutput - Log ra file hay không
	FileOutput bool

	// FilePath - Đường dẫn file log
	FilePath string

	// JSONFormat - Dùng JSON format hay text format
	JSONFormat bool

	// MaxFileSize - Kích thước tối đa của file log (MB) trước khi rotate
	MaxFileSize int

	// MaxBackups - Số lượng file backup giữ lại
	MaxBackups int

	// MaxAge - Số ngày giữ file log cũ
	MaxAge int

	// LogLevel - Level tối thiểu để log (debug, info, warn, error)
	LogLevel string
}

// DefaultLoggerOptions trả về cấu hình mặc định
func DefaultLoggerOptions() LoggerOptions {
	return LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    true,
		FilePath:      "logs/errors.log",
		JSONFormat:    true,
		MaxFileSize:   10,
		MaxBackups:    5,
		MaxAge:        30,
		LogLevel:      "error",
	}
}

// InitLogger khởi tạo logger với custom options
//
// Example:
//
//	config.InitLogger(config.LoggerOptions{
//	    ConsoleOutput: true,
//	    FileOutput: true,
//	    FilePath: "logs/app.log",
//	    JSONFormat: true,
//	})
func InitLogger(opts LoggerOptions) {
	logger := logrus.New()

	// Cấu hình output destinations
	var writers []io.Writer

	// Console output
	if opts.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}

	// File output với rotation
	if opts.FileOutput {
		// Tạo thư mục logs nếu chưa có
		if err := os.MkdirAll("logs", 0755); err != nil {
			logrus.Errorf("Cannot create logs directory: %v", err)
		}

		logFile := &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    opts.MaxFileSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		}
		writers = append(writers, logFile)
	}

	// Set output
	if len(writers) > 0 {
		logger.SetOutput(io.MultiWriter(writers...))
	}

	// Cấu hình formatter
	if opts.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set log level
	level, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logger.SetLevel(level)

	// Wrap và set vào core
	logrusLogger := &LogrusLogger{logger: logger}
	core.SetLogger(logrusLogger)

	logger.Info("✓ GoErrorKit logger initialized")
}

// InitDefaultLogger khởi tạo logger với cấu hình mặc định
//
// Example:
//
//	config.InitDefaultLogger()
func InitDefaultLogger() {
	InitLogger(DefaultLoggerOptions())
}
