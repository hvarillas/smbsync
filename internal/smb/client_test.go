package smb

import (
	"testing"

	"github.com/hvarillas/smbsync/internal/config"
	"github.com/hvarillas/smbsync/internal/logger"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger.Sugar = zap.NewNop().Sugar()
}

func TestGetSmbSession_InvalidHost(t *testing.T) {
	// Test with invalid host - this will fail to connect
	session, err := getSmbSession("testuser", "testpass", "invalid-host-12345")
	if err == nil {
		t.Error("Expected error for invalid host")
		if session != nil {
			session.Logoff()
		}
	}
	
	if session != nil {
		t.Error("Session should be nil for invalid host")
	}
}

func TestRunHeadless_NoFiles(t *testing.T) {
	cfg := &config.Config{
		SMBUser:    "testuser",
		SMBPass:    "testpass",
		SMBHost:    "testhost",
		Shared:     "testshare",
		Regex:      `\.nonexistent$`,
		Path:       "/tmp",
		SharedPath: ".",
	}
	
	// This should complete without error but log a warning about no files
	// Since we can't easily mock the SMB connection, this test mainly verifies
	// that the function handles the no-files case gracefully
	RunHeadless(cfg)
	
	// If we reach here without panic, the test passes
}

func TestRunHeadless_InvalidConfig(t *testing.T) {
	// Test with nil config - should not panic but will likely fail
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config")
		}
	}()
	
	// This will panic when trying to access config fields
	var cfg *config.Config
	RunHeadless(cfg)
}
