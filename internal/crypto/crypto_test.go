package crypto

import (
	"os"
	"testing"
)

func TestEncryptDecryptString(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       string
	}{
		{
			name:      "basic text",
			plaintext: "hello world",
			key:       "1234567890123456",
		},
		{
			name:      "empty string",
			plaintext: "",
			key:       "1234567890123456",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			key:       "abcdefghijklmnop",
		},
		{
			name:      "unicode text",
			plaintext: "こんにちは世界",
			key:       "1234567890123456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetEncryptionKey(tc.key)
			
			encrypted, err := EncryptString(tc.plaintext)
			if err != nil {
				t.Fatalf("EncryptString() error = %v", err)
			}
			
			if encrypted == tc.plaintext {
				t.Error("Encrypted text should not equal plaintext")
			}
			
			decrypted, err := DecryptString(encrypted)
			if err != nil {
				t.Fatalf("DecryptString() error = %v", err)
			}
			
			if decrypted != tc.plaintext {
				t.Errorf("Expected decrypted text '%s', got '%s'", tc.plaintext, decrypted)
			}
		})
	}
}

func TestEncryptDecryptPassword(t *testing.T) {
	SetEncryptionKey("1234567890123456")
	password := "mySecretPassword123"
	
	encrypted, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}
	
	if encrypted == password {
		t.Error("Encrypted password should not equal plaintext password")
	}
	
	decrypted, err := DecryptPassword(encrypted)
	if err != nil {
		t.Fatalf("DecryptPassword() error = %v", err)
	}
	
	if decrypted != password {
		t.Errorf("Expected decrypted password '%s', got '%s'", password, decrypted)
	}
}

func TestInvalidKeyLength(t *testing.T) {
	testCases := []string{
		"short",
		"toolongkey1234567890",
	}
	
	for _, key := range testCases {
		t.Run("key_length_"+key, func(t *testing.T) {
			SetEncryptionKey(key)
			
			_, err := EncryptString("test")
			if err == nil {
				t.Error("Expected error for invalid key length")
			}
		})
	}
	
	// Test empty key separately - it should use default key
	t.Run("empty_key_uses_default", func(t *testing.T) {
		SetEncryptionKey("")
		
		_, err := EncryptString("test")
		if err != nil {
			t.Errorf("Empty key should use default key, got error: %v", err)
		}
	})
}

func TestGetEncryptionKey(t *testing.T) {
	// Test with set key
	testKey := "1234567890123456"
	SetEncryptionKey(testKey)
	
	key := getEncryptionKey()
	if key != testKey {
		t.Errorf("Expected key '%s', got '%s'", testKey, key)
	}
	
	// Test with environment variable
	SetEncryptionKey("")
	os.Setenv("ENCRYPTION_KEY", "envkey1234567890")
	
	key = getEncryptionKey()
	if key != "envkey1234567890" {
		t.Errorf("Expected key from env 'envkey1234567890', got '%s'", key)
	}
	
	// Test with default key
	os.Unsetenv("ENCRYPTION_KEY")
	SetEncryptionKey("")
	
	key = getEncryptionKey()
	if key != "0ED30B7FFA59AFE9" {
		t.Errorf("Expected default key '0ED30B7FFA59AFE9', got '%s'", key)
	}
}

func TestDecryptInvalidData(t *testing.T) {
	SetEncryptionKey("1234567890123456")
	
	testCases := []struct {
		name string
		data string
	}{
		{
			name: "invalid base64",
			data: "invalid-base64!@#",
		},
		{
			name: "too short data",
			data: "dGVzdA==", // "test" in base64, too short for nonce
		},
		{
			name: "corrupted data",
			data: "YWJjZGVmZ2hpams=", // valid base64 but invalid encrypted data
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptString(tc.data)
			if err == nil {
				t.Error("Expected error for invalid encrypted data")
			}
		})
	}
}

func TestSetEncryptionKey(t *testing.T) {
	originalKey := encryptionKey
	defer func() { encryptionKey = originalKey }()
	
	testKey := "newtestkey1234567"
	SetEncryptionKey(testKey)
	
	if encryptionKey != testKey {
		t.Errorf("Expected encryption key '%s', got '%s'", testKey, encryptionKey)
	}
}
