package middleware

import (
	"context"
	//	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Oriseer/workout_tracker/api"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = os.Getenv("JWT_KEY")

type username string

func JwtAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		authParts := strings.Split(authHeader, " ")
		if authParts[0] != "Bearer" || len(authParts) != 2 {
			api.StatusBadRequestServerError(w, api.ErrInvalidToken)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(authParts[1], func(t *jwt.Token) (any, error) {
			// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			// }

			return []byte(jwtKey), nil
		})

		if err != nil {
			api.StatusBadRequestServerError(w, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			api.StatusBadRequestServerError(w, api.ErrInvalidExpredToken)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if mapClaim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var username username
			ctx := context.WithValue(r.Context(), username, mapClaim["username"])
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			api.StatusBadRequestServerError(w, api.ErrInvalidTokenClaims)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
