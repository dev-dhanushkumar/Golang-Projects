package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	level       LogLevel
	output      io.Writer
	mu          sync.Mutex
	enableJSON  bool
	enableColor bool
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// InitLogger initializes the global logger
func InitLogger(level string, enableJSON bool) {
	once.Do(func() {
		logLevel := ParseLogLevel(level)
		defaultLogger = &Logger{
			level:       logLevel,
			output:      os.Stdout,
			enableJSON:  enableJSON,
			enableColor: !enableJSON, // Disable colors in JSON mode
		}
	})
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		InitLogger("INFO", false)
	}
	return defaultLogger
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if level < l.level {
		return // Don't log if below current level
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get caller info (file and line)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		// Get only the filename, not full path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
	}

	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:     level.String(),
		Message:   message,
		File:      file,
		Line:      line,
		Fields:    fields,
	}

	if l.enableJSON {
		l.writeJSON(entry)
	} else {
		l.writeText(entry, level)
	}
}

// writeJSON writes log in JSON format
func (l *Logger) writeJSON(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}
	fmt.Fprintln(l.output, string(data))
}

// writeText writes log in human-readable format
func (l *Logger) writeText(entry LogEntry, level LogLevel) {
	var output strings.Builder

	// Add color if enabled
	if l.enableColor {
		output.WriteString(level.Color())
	}

	// Format: [2024-01-01 12:00:00] INFO [file.go:42] Message
	output.WriteString(fmt.Sprintf("[%s] %-5s", entry.Timestamp, entry.Level))

	if entry.File != "" {
		output.WriteString(fmt.Sprintf(" [%s:%d]", entry.File, entry.Line))
	}

	output.WriteString(fmt.Sprintf(" %s", entry.Message))

	// Add fields if present
	if len(entry.Fields) > 0 {
		output.WriteString(" |")
		for key, value := range entry.Fields {
			output.WriteString(fmt.Sprintf(" %s=%v", key, value))
		}
	}

	// Reset color if enabled
	if l.enableColor {
		output.WriteString(ResetColor)
	}

	fmt.Fprintln(l.output, output.String())
}

// Debug logs debug level message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, message, f)
}

// Info logs info level message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, message, f)
}

// Warn logs warning level message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, message, f)
}

// Error logs error level message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ERROR, message, f)
}

// Fatal logs fatal level message and exits
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, message, f)
	os.Exit(1)
}

// With creates a new logger with predefined fields (for context)
func (l *Logger) With(fields map[string]interface{}) *ContextLogger {
	return &ContextLogger{
		logger: l,
		fields: fields,
	}
}

// ContextLogger wraps Logger with predefined fields
type ContextLogger struct {
	logger *Logger
	fields map[string]interface{}
}

func (cl *ContextLogger) Debug(message string, additionalFields ...map[string]interface{}) {
	fields := cl.mergeFields(additionalFields)
	cl.logger.Debug(message, fields)
}

func (cl *ContextLogger) Info(message string, additionalFields ...map[string]interface{}) {
	fields := cl.mergeFields(additionalFields)
	cl.logger.Info(message, fields)
}

func (cl *ContextLogger) Warn(message string, additionalFields ...map[string]interface{}) {
	fields := cl.mergeFields(additionalFields)
	cl.logger.Warn(message, fields)
}

func (cl *ContextLogger) Error(message string, additionalFields ...map[string]interface{}) {
	fields := cl.mergeFields(additionalFields)
	cl.logger.Error(message, fields)
}

func (cl *ContextLogger) Fatal(message string, additionalFields ...map[string]interface{}) {
	fields := cl.mergeFields(additionalFields)
	cl.logger.Fatal(message, fields)
}

func (cl *ContextLogger) mergeFields(additionalFields []map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Copy context fields
	for k, v := range cl.fields {
		merged[k] = v
	}

	// Add additional fields
	if len(additionalFields) > 0 {
		for k, v := range additionalFields[0] {
			merged[k] = v
		}
	}

	return merged
}

// Global helper functions
func Debug(message string, fields ...map[string]interface{}) {
	GetLogger().Debug(message, fields...)
}

func Info(message string, fields ...map[string]interface{}) {
	GetLogger().Info(message, fields...)
}

func Warn(message string, fields ...map[string]interface{}) {
	GetLogger().Warn(message, fields...)
}

func Error(message string, fields ...map[string]interface{}) {
	GetLogger().Error(message, fields...)
}

func Fatal(message string, fields ...map[string]interface{}) {
	GetLogger().Fatal(message, fields...)
}
