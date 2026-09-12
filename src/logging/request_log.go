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

// func RequestLog() {
// 	slog.Info(
// 		"request started",
// 		"request_id", requestlog.GetRequestID(r.Context()),
// 		"user_id", auth.User.UserID,
// 		"band_id", auth.CurrentBand.BandID,
// 		"method", r.Method,
// 		"path", r.URL.Path,
// 	)
// }
