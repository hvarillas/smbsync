package smb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hvarillas/smbsync/internal/logger"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger.Sugar = zap.NewNop().Sugar()
}

func TestGetRegexFiles(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "smbsync_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files
	testFiles := []string{
		"test1.txt",
		"test2.log",
		"backup.bak",
		"document.pdf",
		"image.jpg",
		"script.sh",
	}
	
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}
	
	// Create a subdirectory (should be ignored)
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	testCases := []struct {
		name     string
		regex    string
		expected []string
	}{
		{
			name:     "match txt files",
			regex:    `\.txt$`,
			expected: []string{"test1.txt"},
		},
		{
			name:     "match log and bak files",
			regex:    `\.(log|bak)$`,
			expected: []string{"test2.log", "backup.bak"},
		},
		{
			name:     "match all files starting with test",
			regex:    `^test`,
			expected: []string{"test1.txt", "test2.log"},
		},
		{
			name:     "match no files",
			regex:    `\.xyz$`,
			expected: []string{},
		},
		{
			name:     "match all files",
			regex:    `.*`,
			expected: testFiles,
		},
		{
			name:     "case insensitive match",
			regex:    `\.PDF$`,
			expected: []string{"document.pdf"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getRegexFiles(tc.regex, tempDir)
			
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d files, got %d", len(tc.expected), len(result))
				t.Errorf("Expected: %v", tc.expected)
				t.Errorf("Got: %v", result)
				return
			}
			
			// Convert to map for easier comparison
			resultMap := make(map[string]bool)
			for _, file := range result {
				resultMap[file] = true
			}
			
			for _, expected := range tc.expected {
				if !resultMap[expected] {
					t.Errorf("Expected file %s not found in result", expected)
				}
			}
		})
	}
}

func TestGetRegexFiles_InvalidRegex(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "smbsync_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test with invalid regex
	result := getRegexFiles("[invalid", tempDir)
	if result != nil {
		t.Error("Expected nil result for invalid regex")
	}
}

func TestGetRegexFiles_NonexistentDirectory(t *testing.T) {
	result := getRegexFiles(`.*`, "/nonexistent/directory")
	if result != nil {
		t.Error("Expected nil result for nonexistent directory")
	}
}

func TestGetRegexFiles_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "smbsync_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	result := getRegexFiles(`.*`, tempDir)
	if len(result) != 0 {
		t.Errorf("Expected empty result for empty directory, got %v", result)
	}
}
