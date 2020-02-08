package app

import (
	"context"
	. "goals/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func (s *server) withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)

		if (*r).Method == "OPTIONS" {
			return
		}
		h(w, r)
	}
}

func (s *server) privateRoute(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitTokenHeader := strings.Split(r.Header.Get("Authorization"), " ")

		if len(splitTokenHeader) != 2 || splitTokenHeader[0] != "Bearer" || splitTokenHeader[1] == "" {
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}

		tokenString := splitTokenHeader[1]

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return s.tokenSigningKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				logrus.WithFields(logrus.Fields{
					"package":     "app",
					"function":    "jwt.ParseWithClaims",
					"error":       err,
					"tokenString": tokenString,
				}).Warning("Token invalid signature")

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

		var user User
		s.db.First(&user, "username = ?", claims.Username)
		ctx := context.WithValue(r.Context(), "User", user)
		r.WithContext(ctx)
		h(w, r.WithContext(ctx))
	}
}
