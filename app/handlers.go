package app

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func (s *server) respond(w http.ResponseWriter, data interface{}, status int) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(data)
		_, err := w.Write(payload)
		if err != nil {
			status = http.StatusInternalServerError
		}
	}

	if status != 200 {
		w.WriteHeader(status)
	}
}

func (s *server) decode(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var requestData request
		err := s.decode(r, &requestData)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		tokenFactory := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": requestData.Username,
			"password": requestData.Password,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := tokenFactory.SignedString(s.signingKey)

		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		resp := response{Token: tokenString}
		s.respond(w, resp, http.StatusOK)
	}
}
