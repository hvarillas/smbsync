package config

import (
	"testing"

	"github.com/hvarillas/smbsync/internal/crypto"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config with plain password",
			config: &Config{
				SMBUser: "testuser",
				SMBPass: "testpass",
				SMBHost: "testhost",
				Shared:  "testshare",
			},
			wantErr: false,
		},
		{
			name: "missing user",
			config: &Config{
				SMBPass: "testpass",
				SMBHost: "testhost",
				Shared:  "testshare",
			},
			wantErr: true,
		},
		{
			name: "missing password and encrypted password",
			config: &Config{
				SMBUser: "testuser",
				SMBHost: "testhost",
				Shared:  "testshare",
			},
			wantErr: true,
		},
		{
			name: "missing host",
			config: &Config{
				SMBUser: "testuser",
				SMBPass: "testpass",
				Shared:  "testshare",
			},
			wantErr: true,
		},
		{
			name: "missing shared",
			config: &Config{
				SMBUser: "testuser",
				SMBPass: "testpass",
				SMBHost: "testhost",
			},
			wantErr: true,
		},
		{
			name: "invalid encryption key length",
			config: &Config{
				SMBUser:       "testuser",
				SMBPass:       "testpass",
				SMBHost:       "testhost",
				Shared:        "testshare",
				EncryptionKey: "short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_EncryptPassword(t *testing.T) {
	config := &Config{SMBPass: "testpassword"}
	
	encrypted, err := config.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}
	
	if encrypted == "" {
		t.Error("EncryptPassword() returned empty string")
	}
	
	if encrypted == "testpassword" {
		t.Error("EncryptPassword() returned plain text password")
	}
}

func TestConfig_EncryptString(t *testing.T) {
	config := &Config{}
	testText := "sensitive data"
	
	encrypted, err := config.EncryptString(testText)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	
	if encrypted == "" {
		t.Error("EncryptString() returned empty string")
	}
	
	if encrypted == testText {
		t.Error("EncryptString() returned plain text")
	}
}

func TestConfig_ValidateWithEncryptedPassword(t *testing.T) {
	// Set up encryption key
	crypto.SetEncryptionKey("1234567890123456")
	
	// Encrypt a password
	encrypted, err := crypto.EncryptPassword("testpass")
	if err != nil {
		t.Fatalf("Failed to encrypt password: %v", err)
	}
	
	config := &Config{
		SMBUser:       "testuser",
		SMBHost:       "testhost",
		Shared:        "testshare",
		EncryptedPass: encrypted,
	}
	
	err = config.Validate()
	if err != nil {
		t.Errorf("Config.Validate() with encrypted password error = %v", err)
	}
	
	if config.SMBPass != "testpass" {
		t.Errorf("Expected decrypted password 'testpass', got '%s'", config.SMBPass)
	}
}

func TestLoad(t *testing.T) {
	// Set some test values
	smbUser = "testuser"
	smbPass = "testpass"
	smbHost = "testhost"
	shared = "testshare"
	
	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	
	if config.SMBUser != "testuser" {
		t.Errorf("Expected SMBUser 'testuser', got '%s'", config.SMBUser)
	}
	
	if config.SMBPass != "testpass" {
		t.Errorf("Expected SMBPass 'testpass', got '%s'", config.SMBPass)
	}
	
	if config.SMBHost != "testhost" {
		t.Errorf("Expected SMBHost 'testhost', got '%s'", config.SMBHost)
	}
	
	if config.Shared != "testshare" {
		t.Errorf("Expected Shared 'testshare', got '%s'", config.Shared)
	}
}
