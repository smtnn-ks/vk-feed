package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware(d deps, next func(w http.ResponseWriter, r *http.Request), isOpt bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if authParts[0] != "Bearer" {
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(authParts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(d.jwtSecret), nil
		}, jwt.WithExpirationRequired())
		if err != nil {
			log.Println(err)
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId := claims["sub"].(string)
		expiresAt := int64(claims["exp"].(float64))
		if expiresAt < time.Now().UTC().Unix() {
			if isOpt {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.Header.Add("userid", userId)
		next(w, r)
	}
}
