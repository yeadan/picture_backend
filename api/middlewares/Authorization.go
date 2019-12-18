package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/yeadan/picture_backend/api/data"
	"github.com/yeadan/picture_backend/lib"
)

var UserKey = "current_user"

// AuthUser - Autorización con bearer token. No lo evalúa si es un signup o un login
func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL != nil && (r.URL.Path == "/" || r.URL.Path == "/users/login" || (r.Method == "POST" && r.URL.Path == "/users")) {
			next.ServeHTTP(w, r)
			return
		}
		auto := r.Header.Get("Authorization")
		if len(auto) > 0 && strings.Contains(auto, "Bearer ") {
			tokenString := strings.Split(auto, " ")[1]
			userValid := lib.GetUserTokenCache(tokenString, data.GetCacheClient())
			if userValid != nil {
				ctx := context.WithValue(r.Context(), UserKey, userValid)
				newReq := r.WithContext(ctx)
				next.ServeHTTP(w, newReq)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}
