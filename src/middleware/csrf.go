package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

const csrfTokenCookieName = "csrf_token"

func CSRFTokenCookie(next http.Handler, secure bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     csrfTokenCookieName,
			Value:    csrf.Token(r),
			Path:     "/",
			Secure:   secure,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		})

		next.ServeHTTP(w, r)
	})
}

func CSRFPlaintext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
	})
}
