package server

import (
	"encoding/json"
	. "goals/app/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) respond(w http.ResponseWriter, data interface{}, status int) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(data)
		_, err := w.Write(payload)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"function": "ResponseWriter.Write",
				"error":    err,
				"payload":  payload,
			}).Warning("Internal server error")

			status = http.StatusInternalServerError
		}
	}

	if status != 200 {
		w.WriteHeader(status)
	}
}

func (s *Server) createToken(username string) (string, error) {
	tokenFactory := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	return tokenFactory.SignedString(s.tokenSigningKey)
}

func (s *Server) handleRegister() http.HandlerFunc {
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
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleRegister",
				"function": "json.NewDecoder",
				"error":    err,
				"data":     r.Body,
			}).Warning("Failed to decode request body")

			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		var newUser User
		s.db.First(&newUser, "username = ?", requestData.Username)

		if newUser.Username != "" {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		s.db.First(&newUser, "email = ?", requestData.Email)

		if newUser.Email != "" {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), 8)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleRegister",
				"function": "bcrypt.GenerateFromPassword",
				"error":    err,
			}).Warning("Failed to generate password hash")

			s.respond(w, nil, http.StatusInternalServerError)
			return
		}
		newUser = User{
			Username:  requestData.Username,
			Password:  string(hashedPassword),
			Email:     requestData.Email,
			FirstName: requestData.FirstName,
			LastName:  requestData.LastName,
		}

		s.db.Create(&newUser)

		tokenString, err := s.createToken(requestData.Username)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleRegister",
				"function": "server.createToken",
				"error":    err,
				"data":     requestData.Username,
			}).Warning("Failed to create token")

			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		resp := response{Token: tokenString}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
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
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleLogin",
				"function": "json.NewDecoder",
				"error":    err,
				"data":     r.Body,
			}).Warning("Failed to decode request body")

			s.respond(w, nil, http.StatusBadRequest)
			return
		}
		var user User
		s.db.First(&user, "username = ?", requestData.Username)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))

		if user.Username == "" || err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleLogin",
				"function": "bcrypt.CompareHashAndPassword",
				"error":    err,
			}).Warning("Failed to login")

			s.respond(w, nil, http.StatusUnauthorized)
			return
		}

		tokenString, err := s.createToken(requestData.Username)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleLogin",
				"function": "server.createToken",
				"error":    err,
				"data":     requestData.Username,
			}).Warning("Failed to create token")

			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		resp := response{Token: tokenString}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *Server) handleWhoAmI() http.HandlerFunc {
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
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleWhoAmI",
				"function": "copier.Copy",
				"error":    err,
			}).Warning("Failed to write response")

			s.respond(w, nil, http.StatusInternalServerError)
			return
		}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	type response struct {
		Result string `json:"result"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{Result: "pong"}
		s.respond(w, resp, http.StatusOK)
	}
}

func (s *Server) handleGetAreas() http.HandlerFunc {
	type areasResponse struct {
		Id          int    `json:"id"`
		Weight      int    `json:"weight"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		IsFavourite bool   `json:"is_favourite"`
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

func (s *Server) handleCreateAreas() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		IsFavourite bool   `json:"is_favourite"`
		Weight      int    `json:"weight"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var requestData request
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package":  "app",
				"handler":  "handleCreateAreas",
				"function": "json.NewDecoder",
				"error":    err,
				"data":     r.Body,
			}).Warning("Failed to decode request body")
		}

		user := r.Context().Value("User").(User)

		newArea := Area{
			Name:        requestData.Name,
			Description: requestData.Description,
			Icon:        requestData.Icon,
			IsFavourite: requestData.IsFavourite,
			Weight:      requestData.Weight,
			UserID:      user.ID,
		}

		s.db.Create(&newArea)

		s.respond(w, nil, http.StatusCreated)
	}
}
