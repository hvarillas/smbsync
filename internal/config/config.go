package config

import (
	"fmt"

	"github.com/hvarillas/smbsync/internal/crypto"
	"github.com/spf13/cobra"
)

type Config struct {
	SMBUser        string
	SMBPass        string
	SMBHost        string
	Regex          string
	Path           string
	Shared         string
	SharedPath     string
	DeleteAfter    bool
	Zippy          bool
	LogPath        string
	LogLevel       string
	EncryptedPass  string
	GenerateCrypto bool
	EncryptText    string
	EncryptionKey  string
}

func Load() (*Config, error) {
	return &Config{
		SMBUser:        smbUser,
		SMBPass:        smbPass,
		SMBHost:        smbHost,
		Regex:          regex,
		Path:           path,
		Shared:         shared,
		SharedPath:     sharedPath,
		DeleteAfter:    deleteAfter,
		Zippy:          zippy,
		LogPath:        logPath,
		LogLevel:       logLevel,
		EncryptedPass:  encryptedPass,
		GenerateCrypto: generateCrypto,
		EncryptText:    encryptText,
		EncryptionKey:  encryptionKey,
	}, nil
}

func (c *Config) Validate() error {
	if c.EncryptionKey != "" {
		if len(c.EncryptionKey) != 16 {
			return fmt.Errorf("la clave de encriptación debe tener exactamente 16 bytes")
		}
		crypto.SetEncryptionKey(c.EncryptionKey)
	}

	if c.SMBUser == "" || (c.SMBPass == "" && c.EncryptedPass == "") || c.SMBHost == "" || c.Shared == "" {
		return fmt.Errorf("user, (pass o encrypted-pass), host, y shared son requeridos")
	}

	if c.EncryptedPass != "" {
		decrypted, err := crypto.DecryptPassword(c.EncryptedPass)
		if err != nil {
			return fmt.Errorf("error al desencriptar contraseña: %w", err)
		}
		c.SMBPass = decrypted
	}

	return nil
}

func (c *Config) EncryptPassword() (string, error) {
	return crypto.EncryptPassword(c.SMBPass)
}

func (c *Config) EncryptString(text string) (string, error) {
	return crypto.EncryptString(text)
}

var (
	smbUser        string
	smbPass        string
	smbHost        string
	regex          string
	path           string
	shared         string
	sharedPath     string
	deleteAfter    bool
	zippy          bool
	logPath        string
	logLevel       string
	encryptedPass  string
	generateCrypto bool
	encryptText    string
	encryptionKey  string
)

func InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&smbUser, "user", "u", "", "SMB user name (required)")
	cmd.PersistentFlags().StringVarP(&smbPass, "pass", "p", "", "SMB password (required if encrypted-pass not provided)")
	cmd.PersistentFlags().StringVar(&smbHost, "host", "", "SMB host (required)")
	cmd.PersistentFlags().StringVar(&encryptedPass, "encrypted-pass", "", "Encrypted SMB password (alternative to --pass)")
	cmd.PersistentFlags().BoolVar(&generateCrypto, "generate-encrypted", false, "Generate encrypted password from --pass flag")
	cmd.PersistentFlags().StringVar(&encryptText, "encrypt-text", "", "Encrypt any string using AES-GCM")
	cmd.PersistentFlags().StringVar(&encryptionKey, "encryption-key", "", "16-byte encryption key for AES (overrides ENCRYPTION_KEY env var)")
	cmd.PersistentFlags().StringVarP(&shared, "shared", "s", "", "Shared resource SMB (required)")
	cmd.PersistentFlags().StringVarP(&regex, "regex", "r", "", "Regex to filter local files")
	cmd.PersistentFlags().StringVarP(&path, "path", "", ".", "Base path for local files to copy")
	cmd.PersistentFlags().StringVarP(&sharedPath, "sharedPath", "", ".", "Relative destination path on the share")
	cmd.PersistentFlags().BoolVarP(&deleteAfter, "delete", "d", false, "Delete local file after successful verification")
	cmd.PersistentFlags().BoolVarP(&zippy, "zip", "z", false, "Compress files before copying")
	cmd.PersistentFlags().StringVarP(&logPath, "log", "l", "smbsync.log", "Path to the log file")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set the log level (debug, info, warn, error)")
}
