package middleware

import (
	"bandplan/src/handlers"
	requestlog "bandplan/src/logging"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, band, err := handlers.HelperGetAuthenticatedUserAndBand(r)
		if err != nil {
			slog.Error(
				"Unable to authorize user",
				"requets_id", requestlog.GetRequestID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"hx_header", r.Header.Get("HX-Request"),
			)
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		auth := handlers.AuthContext{
			User:        user,
			CurrentBand: band,
		}

		ctx := context.WithValue(
			r.Context(),
			handlers.AuthContextKey,
			auth,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", "error", err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()

		ctx := context.WithValue(
			r.Context(),
			requestlog.RequestIDKey,
			requestID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		requestID := requestlog.GetRequestID(r.Context())

		auth, ok := r.Context().Value(handlers.AuthContextKey).(handlers.AuthContext)

		if ok {
			slog.Info(
				r.Method+" "+r.URL.Path+"  started - ",
				"request_id", requestID[:13],
				"user_id", auth.User.UserID,
				"band_id", auth.CurrentBand.BandID[:8],
			)
		}

		next.ServeHTTP(w, r)

		duration := time.Since(start).Round(100 * time.Microsecond)

		slog.Info(
			r.Method+" "+r.URL.Path+"  completed - ",
			"request_id", requestID[:13],
			"duration", duration,
		)
	})
}
