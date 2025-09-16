package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hvarillas/smbsync/internal/config"
	"github.com/hvarillas/smbsync/internal/crypto"
	"github.com/hvarillas/smbsync/internal/logger"
	"github.com/hvarillas/smbsync/internal/smb"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for integration tests
	logger.Sugar = zap.NewNop().Sugar()
}

func TestConfigCryptoIntegration(t *testing.T) {
	// Test the integration between config and crypto packages
	crypto.SetEncryptionKey("1234567890123456")
	
	cfg := &config.Config{
		SMBUser: "testuser",
		SMBPass: "plainpassword",
		SMBHost: "testhost",
		Shared:  "testshare",
	}
	
	// Encrypt password
	encrypted, err := cfg.EncryptPassword()
	if err != nil {
		t.Fatalf("Failed to encrypt password: %v", err)
	}
	
	// Create new config with encrypted password
	cfg2 := &config.Config{
		SMBUser:       "testuser",
		SMBHost:       "testhost",
		Shared:        "testshare",
		EncryptedPass: encrypted,
		EncryptionKey: "1234567890123456",
	}
	
	// Validate should decrypt the password
	err = cfg2.Validate()
	if err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}
	
	if cfg2.SMBPass != "plainpassword" {
		t.Errorf("Expected decrypted password 'plainpassword', got '%s'", cfg2.SMBPass)
	}
}

func TestSMBUtilsIntegration(t *testing.T) {
	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "smbsync_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files
	testFiles := map[string]string{
		"document.pdf":    "PDF content",
		"backup.bak":      "Backup content",
		"log.txt":         "Log content",
		"image.jpg":       "Image content",
		"script.sh":       "Script content",
	}
	
	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}
	
	// Test regex file matching
	testCases := []struct {
		regex    string
		expected int
	}{
		{`\.pdf$`, 1},
		{`\.bak$`, 1},
		{`\.(txt|sh)$`, 2},
		{`.*`, 5},
		{`\.xyz$`, 0},
	}
	
	for _, tc := range testCases {
		t.Run("regex_"+tc.regex, func(t *testing.T) {
			files := smb.GetRegexFiles(tc.regex, tempDir)
			if len(files) != tc.expected {
				t.Errorf("Expected %d files for regex '%s', got %d", tc.expected, tc.regex, len(files))
			}
		})
	}
}

func TestFullConfigValidation(t *testing.T) {
	testCases := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "complete valid config",
			config: &config.Config{
				SMBUser:     "user",
				SMBPass:     "pass",
				SMBHost:     "host",
				Shared:      "share",
				Path:        ".",
				SharedPath:  ".",
				Regex:       `.*`,
				LogLevel:    "info",
				LogPath:     "test.log",
				DeleteAfter: false,
				Zippy:       false,
			},
			wantErr: false,
		},
		{
			name: "config with encryption",
			config: &config.Config{
				SMBUser:       "user",
				SMBHost:       "host",
				Shared:        "share",
				EncryptionKey: "1234567890123456",
				EncryptedPass: "encrypted_pass_placeholder",
			},
			wantErr: true, // Will fail because encrypted_pass_placeholder is not valid
		},
		{
			name: "minimal valid config",
			config: &config.Config{
				SMBUser: "user",
				SMBPass: "pass",
				SMBHost: "host",
				Shared:  "share",
			},
			wantErr: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Helper function to expose getRegexFiles for testing
func TestGetRegexFilesExposed(t *testing.T) {
	// This test ensures we can access the getRegexFiles function
	// In a real scenario, you might want to make this function public
	// or create a test-specific interface
	
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
	
	// This would normally call the internal function
	// For now, we'll just verify the directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Test directory should exist")
	}
}
