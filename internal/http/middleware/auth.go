package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/response"
)

type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	cookieName            = "auth"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(cookieName)
			if err != nil || !validCookie(c.Value, secretKey) {
				response.Unauthorized(w)
				return
			}

			parts := strings.Split(c.Value, "|")
			userID := parts[0]

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Sign(value, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func validCookie(value, secretKey string) bool {
	parts := strings.Split(value, "|")
	if len(parts) != 2 {
		return false
	}
	return Sign(parts[0], secretKey) == parts[1]
}
