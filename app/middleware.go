package app

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func (s *server) privateRoute(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return s.tokenSigningKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				s.respond(w, nil, http.StatusUnauthorized)
				return
			}
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}
		if !tkn.Valid {
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}
