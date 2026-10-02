package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
)

func TestCSRFTokenCookieSupportsHeaderValidation(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	handler := csrf.Protect(key, csrf.Secure(false))(
		CSRFTokenCookie(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}), false),
	)
	handler = CSRFPlaintext(handler)

	getResponse := httptest.NewRecorder()
	handler.ServeHTTP(getResponse, httptest.NewRequest(http.MethodGet, "/", nil))

	var tokenCookie *http.Cookie
	for _, cookie := range getResponse.Result().Cookies() {
		if cookie.Name == csrfTokenCookieName {
			tokenCookie = cookie
			break
		}
	}

	if tokenCookie == nil {
		t.Fatal("CSRF token cookie was not set")
	}
	if tokenCookie.Value == "" {
		t.Fatal("CSRF token cookie was empty")
	}
	if tokenCookie.HttpOnly {
		t.Fatal("CSRF token cookie must be readable by HTMX configuration")
	}
	if tokenCookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("CSRF token cookie SameSite = %v, want %v", tokenCookie.SameSite, http.SameSiteLaxMode)
	}

	postRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range getResponse.Result().Cookies() {
		postRequest.AddCookie(cookie)
	}
	postRequest.Header.Set("X-CSRF-Token", tokenCookie.Value)

	postResponse := httptest.NewRecorder()
	handler.ServeHTTP(postResponse, postRequest)

	if postResponse.Code != http.StatusNoContent {
		t.Fatalf("POST status = %d, want %d", postResponse.Code, http.StatusNoContent)
	}
}

func TestCSRFTokenCookieRejectsMissingHeader(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	handler := csrf.Protect(key, csrf.Secure(false))(
		CSRFTokenCookie(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}), false),
	)
	handler = CSRFPlaintext(handler)

	getResponse := httptest.NewRecorder()
	handler.ServeHTTP(getResponse, httptest.NewRequest(http.MethodGet, "/", nil))

	postRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	for _, cookie := range getResponse.Result().Cookies() {
		postRequest.AddCookie(cookie)
	}

	postResponse := httptest.NewRecorder()
	handler.ServeHTTP(postResponse, postRequest)

	if postResponse.Code != http.StatusForbidden {
		t.Fatalf("POST status = %d, want %d", postResponse.Code, http.StatusForbidden)
	}
}
