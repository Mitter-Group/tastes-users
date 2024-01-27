package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func EncryptAPIKey(apiKey string, secretKey string) (string, error) {
	// Convertir la clave API y la clave secreta a bytes
	apiKeyBytes := []byte(apiKey)
	secretKeyBytes := []byte(secretKey)

	// Crear un bloque de cifrado AES
	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return "", err
	}

	// Generar un nonce aleatorio
	nonce := make([]byte, 12) // Cambiar aes.BlockSize a 12
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Crear un cifrador de bloque con modo Galois/Counter Mode (GCM)
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encriptar la clave API
	encryptedAPIKey := aesgcm.Seal(nil, nonce, apiKeyBytes, nil)

	// Concatenar el nonce y la clave encriptada
	result := append(nonce, encryptedAPIKey...)

	// Devolver la clave encriptada como string base64
	return base64.StdEncoding.EncodeToString(result), nil
}

func DecryptAndValidateAPIKey(encryptedAPIKey string, secretKey string) (string, error) {
	// Decodificar la clave encriptada desde base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedAPIKey)
	if err != nil {
		return "", err
	}

	// Convertir la clave secreta a bytes
	secretKeyBytes := []byte(secretKey)

	// Crear un bloque de cifrado AES
	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return "", err
	}

	// Obtener el tamaño del nonce
	nonceSize := 12 // Cambiar aes.BlockSize a 12

	// Extraer el nonce de los primeros nonceSize bytes de la clave encriptada
	nonce := encryptedBytes[:nonceSize]

	// Resto de la clave encriptada
	encryptedAPIKeyBytes := encryptedBytes[nonceSize:]

	// Crear un cifrador de bloque con modo Galois/Counter Mode (GCM)
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Desencriptar la clave API
	apiKeyBytes, err := aesgcm.Open(nil, nonce, encryptedAPIKeyBytes, nil)
	if err != nil {
		return "", err
	}

	// Devolver la clave API como string
	return string(apiKeyBytes), nil
}
