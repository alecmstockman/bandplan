package middleware

import (
	"bandplan/src/handlers"
	"context"
	"log"
	"net/http"
)

type contextKey string

const authenticatedUserKey contextKey = "authenticated-user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("- Middleware RequireAuth")
		user, band, err := handlers.HelperGetAuthenticatedUserAndBand(r)
		if err != nil {
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
