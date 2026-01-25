package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
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
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			parts := strings.Split(c.Value, "|")
			ctx := context.WithValue(r.Context(), UserIDKey, parts[0])

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

func UserID(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}
