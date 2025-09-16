package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	// Create temporary log file
	tempDir, err := os.MkdirTemp("", "smbsync_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	logPath := filepath.Join(tempDir, "test.log")
	
	testCases := []struct {
		name     string
		logLevel string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"invalid level defaults to info", "invalid"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Remove log file if exists
			os.Remove(logPath)
			
			Init(logPath, tc.logLevel)
			
			if Sugar == nil {
				t.Error("Sugar logger should be initialized")
			}
			
			// Test logging with info level (should always work for info and above)
			Sugar.Info("test message")
			
			// Check if log file was created
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				t.Error("Log file should be created")
			}
			
			// Wait a bit for file write
			time.Sleep(100 * time.Millisecond)
			
			// Check log file content - only check if file exists and has content
			content, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			
			// For warn and error levels, info messages might not be logged
			// So we'll just check that the file was created and logger initialized
			if len(content) == 0 && (tc.logLevel == "debug" || tc.logLevel == "info" || tc.logLevel == "invalid") {
				t.Error("Log file should contain content for debug/info levels")
			}
		})
	}
}

func TestInit_InvalidLogPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid log path")
		}
	}()
	
	// Try to create log file in non-existent directory without permission
	Init("/root/nonexistent/test.log", "info")
}

func TestLogLevels(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "smbsync_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	logPath := filepath.Join(tempDir, "levels.log")
	
	testCases := []struct {
		initLevel string
		logLevel  string
		logFunc   func(string)
		shouldLog bool
	}{
		{"debug", "debug", func(msg string) { Sugar.Debug(msg) }, true},
		{"debug", "info", func(msg string) { Sugar.Info(msg) }, true},
		{"debug", "warn", func(msg string) { Sugar.Warn(msg) }, true},
		{"debug", "error", func(msg string) { Sugar.Error(msg) }, true},
		{"info", "debug", func(msg string) { Sugar.Debug(msg) }, false},
		{"info", "info", func(msg string) { Sugar.Info(msg) }, true},
		{"warn", "info", func(msg string) { Sugar.Info(msg) }, false},
		{"warn", "warn", func(msg string) { Sugar.Warn(msg) }, true},
		{"error", "warn", func(msg string) { Sugar.Warn(msg) }, false},
		{"error", "error", func(msg string) { Sugar.Error(msg) }, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.initLevel+"_"+tc.logLevel, func(t *testing.T) {
			// Remove and recreate log file
			os.Remove(logPath)
			
			Init(logPath, tc.initLevel)
			
			testMessage := "test_" + tc.logLevel + "_message"
			tc.logFunc(testMessage)
			
			// Wait for file write
			time.Sleep(50 * time.Millisecond)
			
			content, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			
			containsMessage := strings.Contains(string(content), testMessage)
			if tc.shouldLog && !containsMessage {
				t.Errorf("Log file should contain message '%s' for level %s", testMessage, tc.initLevel)
			} else if !tc.shouldLog && containsMessage {
				t.Errorf("Log file should not contain message '%s' for level %s", testMessage, tc.initLevel)
			}
		})
	}
}
