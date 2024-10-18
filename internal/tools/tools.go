package tools

import (
	"context"
	"errors"

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

func NewEncryptor(name string, secret []byte) (Encryptor, error) {
	switch name {
	case "aes":
		return NewAESEncryptor(secret)
	case "chacha20-poly1305":
		return NewChaChaEncryptor(secret)
	default:
		return nil, errors.New("unknown encryptor")
	}
}

type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}
