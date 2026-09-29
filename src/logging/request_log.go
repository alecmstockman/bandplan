package requestlog

import (
	"context"
)

type contextKey string

const AuthenticatedUserKey contextKey = "authenticated-user"
const RequestIDKey contextKey = "request-id"

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(RequestIDKey).(string)
	return requestID
}
