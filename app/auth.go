package app

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/glorinli/go-jwt-simple-auth/models"
	u "github.com/glorinli/go-jwt-simple-auth/utils"
	"net/http"
	"os"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Endpoints that need no authentication
		needAuthPaths := []string{"/api/account/me"}
		requestPath := r.URL.Path

		var needAuth = false
		for _, value := range needAuthPaths {
			if value == requestPath {
				needAuth = true
				break
			}
		}

		if !needAuth {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")

		// Token is missing
		if tokenHeader == "" {
			sendInvalidTokenResponse(w, "Missing auth token")
			return
		}

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		fmt.Println("Parse token error:", err)

		if err != nil {
			sendInvalidTokenResponse(w, "Invalid auth token: "+err.Error())
			return
		}

		// Token is invalid
		if !token.Valid {
			sendInvalidTokenResponse(w, "Token is not valid")
			return
		}

		// Auth ok
		fmt.Println("User:", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func sendInvalidTokenResponse(w http.ResponseWriter, message string) {
	response := u.Message(false, message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	u.Respond(w, response)
}
