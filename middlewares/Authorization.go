package middlewares

import (
	"context"
	"fmt"
	"github.com/RGlimadigital/Tareas-Go/data"
	"github.com/RGlimadigital/Tareas-Go/lib"
	"net/http"
	"strings"
)

var UserKey = "current_user"

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Auth")
		if r.URL != nil && (r.URL.Path == "/users/login" || r.URL.Path == "/users") {
			next.ServeHTTP(w, r)
			return
		}
		auto := r.Header.Get("Authorization")
		if len(auto) > 0 && strings.Contains(auto, "Bearer ") {
			tokenString := strings.Split(auto, " ")[1]
			fmt.Println(tokenString)
			userValid := lib.GetUserTokenCache(tokenString, data.GetCacheClient())
			// userValid := lib.GetUserJWT(tokenString)
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
