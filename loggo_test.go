package loggo

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	logger, err := NewLogger(tempFile.Name(), LevelInfo)
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestLogLevels(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	logger, err := NewLogger(tempFile.Name(), LevelDebug)
	require.NoError(t, err)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		message  string
		expected map[string]string
	}{
		{"DEBUG", logger.Debug, "Debug message", map[string]string{"level": "DEBUG", "msg": "Debug message"}},
		{"INFO", logger.Info, "Info message", map[string]string{"level": "INFO", "msg": "Info message"}},
		{"WARN", logger.Warn, "Warn message", map[string]string{"level": "WARN", "msg": "Warn message"}},
		{"ERROR", func(msg string, args ...interface{}) {
			logger.Error(msg, errors.New("test error"))
		}, "Error message", map[string]string{"level": "ERROR", "msg": "Error message", "error": "test error"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(tt.message)

			// Read the log file content
			content, err := os.ReadFile(tempFile.Name())
			require.NoError(t, err)

			// Trim any whitespace
			output := strings.TrimSpace(string(content))

			// Check if the output is JSON
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(output), &logEntry); err == nil {
				// JSON format
				for key, expectedValue := range tt.expected {
					assert.Equal(t, expectedValue, logEntry[key], "Mismatch in field: %s", key)
				}
			} else {
				// Non-JSON format
				for key, expectedValue := range tt.expected {
					assert.Contains(t, output, key+"="+expectedValue, "Missing or incorrect %s", key)
				}
			}

			// Clear the file for the next test
			err = os.Truncate(tempFile.Name(), 0)
			require.NoError(t, err)
		})
	}
}

func TestWithOperation(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	logger, err := NewLogger(tempFile.Name(), LevelInfo)
	require.NoError(t, err)

	operationLogger := logger.WithOperation("test-operation")
	assert.NotNil(t, operationLogger)

	// Capture file output
	operationLogger.Info("Operation log")

	content, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)

	var logEntry map[string]interface{}
	err = json.Unmarshal(content, &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "Operation log", logEntry["msg"])
	assert.Equal(t, "test-operation", logEntry["operationID"])
}

func TestIsDebugEnabled(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name     string
		level    Level
		expected bool
	}{
		{"DebugEnabled", LevelDebug, true},
		{"InfoEnabled", LevelInfo, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tempFile.Name(), tt.level)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, logger.IsDebugEnabled())
		})
	}
}

func TestFatal(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test-log-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	logger, err := NewLogger(tempFile.Name(), LevelInfo)
	require.NoError(t, err)

	// Use a custom exit function to prevent the test from actually exiting
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	logger.Fatal("Fatal error", nil)

	// Read the log file content
	content, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)

	// Convert content to string and trim any whitespace
	output := strings.TrimSpace(string(content))

	// Parse JSON content
	var logEntry map[string]interface{}
	err = json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err, "Failed to parse JSON: %v", err)

	// Check log entry fields
	assert.Equal(t, "ERROR", logEntry["level"], "Incorrect log level")
	assert.Equal(t, "Fatal error", logEntry["msg"], "Incorrect log message")
	assert.Nil(t, logEntry["error"], "Error should be nil")

	// Check exit code
	assert.Equal(t, 1, exitCode, "Incorrect exit code")
}
