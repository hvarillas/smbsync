package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

var encryptionKey string

func SetEncryptionKey(key string) {
	encryptionKey = key
}

func getEncryptionKey() string {
	if encryptionKey != "" {
		return encryptionKey
	}
	if key := os.Getenv("ENCRYPTION_KEY"); key != "" {
		return key
	}
	return "0ED30B7FFA59AFE9"
}

func EncryptString(text string) (string, error) {
	key := getEncryptionKey()
	if len(key) != 16 {
		return "", fmt.Errorf("la clave de encriptación debe tener exactamente 16 bytes")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error al crear cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error al crear GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error al generar nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("error al decodificar base64: %w", err)
	}

	key := getEncryptionKey()
	if len(key) != 16 {
		return "", fmt.Errorf("la clave de encriptación debe tener exactamente 16 bytes")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error al crear cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error al crear GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("texto cifrado muy corto")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error al desencriptar: %w", err)
	}

	return string(plaintext), nil
}

func EncryptPassword(password string) (string, error) {
	return EncryptString(password)
}

func DecryptPassword(encryptedPassword string) (string, error) {
	return DecryptString(encryptedPassword)
}
