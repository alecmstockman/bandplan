package requestlog

import (
	"context"
)

type contextKey string

const AuthenticatedUserKey contextKey = "authenticated-user"
const RequestIDKey contextKey = "request-id"
const ClientIPKey contextKey = "client-ip"

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(RequestIDKey).(string)
	return requestID
}

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ClientIPKey, ip)
}

func GetClientIP(ctx context.Context) string {
	ip, _ := ctx.Value(ClientIPKey).(string)
	return ip
}
