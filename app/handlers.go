package app

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
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

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *server) handleRegister() http.HandlerFunc {
	type request struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	type response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var requestData request
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		var newUser User
		s.db.First(&newUser, "username = ?", requestData.Username)

		if newUser.Username != "" {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		err = copier.Copy(&newUser, &requestData)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		s.db.Create(&newUser)

		tokenFactory := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			Username: requestData.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		})

		tokenString, err := tokenFactory.SignedString(s.tokenSigningKey)

		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		resp := response{Token: tokenString}
		s.respond(w, resp, http.StatusOK)
	}
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
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}
		var user User
		s.db.First(&user, "username = ?", requestData.Username)

		if user.Username == "" || user.Password != requestData.Password {
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}
		tokenFactory := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			Username: requestData.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		})

		tokenString, err := tokenFactory.SignedString(s.tokenSigningKey)

		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		resp := response{Token: tokenString}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	type response struct {
		Username  string `json:"username"`
		ID        uint   `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("User").(User)
		resp := response{}
		err := copier.Copy(&resp, &user)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *server) handlePing() http.HandlerFunc {
	type response struct {
		Result string `json:"result"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{Result: "pong"}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *server) handleGetAreas() http.HandlerFunc {
	type areasResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type response struct {
		Areas []areasResponse `json:"areas"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("User").(User)
		resp := response{Areas: []areasResponse{}}
		s.db.Table("areas").Where("user_id = ?", user.ID).Scan(&resp.Areas)
		s.respond(w, resp, http.StatusOK)
	}
}
