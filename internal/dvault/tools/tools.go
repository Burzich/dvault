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
