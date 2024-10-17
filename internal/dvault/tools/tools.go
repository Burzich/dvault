package tools

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/google/uuid"
)

var RequestID struct{}

func GetRequestIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(RequestID).(string)
	if !ok {
		return ""
	}

	return val
}

func AddXRequestIDToContext(ctx context.Context) context.Context {
	requestID := uuid.NewString()

	return context.WithValue(ctx, RequestID, requestID)
}

func GenerateXRequestID() string {
	return uuid.NewString()
}

func Encrypt(data []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func Decrypt(data []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedData, err := gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
