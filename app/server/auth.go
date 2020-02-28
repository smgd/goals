package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (s *Server) createToken(username string) (string, error) {
	tokenFactory := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	return tokenFactory.SignedString([]byte(s.config.TokenSigningKey))
}
