package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.SugaredLogger, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	// Open log file
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Configure encoder for file (JSON format)
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// Configure encoder for console (human-readable format)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Set log level
	level := zapcore.InfoLevel
	if os.Getenv("ENVIRONMENT") == "development" {
		level = zapcore.DebugLevel
	}

	// Create core that writes to both file and console
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	// Create logger with caller information
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Assign the SugaredLogger to the global variable
	sugar := logger.Sugar()
	return sugar, nil
}
