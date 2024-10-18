package tools

import (
	"context"

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

func NewEncryptor(secret []byte) (Encryptor, error) {
	return nil, nil
}

type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}
