package server

import (
	"encoding/json"
	"goals/app/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
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

	return tokenFactory.SignedString([]byte(s.config.TokenSigningKey))
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

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		user := models.User{
			Username:  requestData.Username,
			Password:  requestData.Password,
			Email:     requestData.Email,
			FirstName: requestData.FirstName,
			LastName:  requestData.LastName,
		}

		if _, err := s.store.User().Create(&user); err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		tokenString, err := s.createToken(requestData.Username)
		if err != nil {
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

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		user, err := s.store.User().FindByUsername(requestData.Username)

		if err != nil || user.ComparePassword(requestData.Password) {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		tokenString, err := s.createToken(requestData.Username)
		if err != nil {
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
		user := r.Context().Value("User").(*models.User)
		resp := response{}
		if err := copier.Copy(&resp, &user); err != nil {
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
		user := r.Context().Value("User").(*models.User)

		userAreas, err := s.store.Area().FindAreasByUserID(user.ID)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
		}

		resp := response{Areas: []areasResponse{}}

		if err := copier.Copy(&resp.Areas, &userAreas); err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

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

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		user := r.Context().Value("User").(*models.User)

		area := models.Area{
			Name:        requestData.Name,
			Description: requestData.Description,
			Icon:        requestData.Icon,
			IsFavourite: requestData.IsFavourite,
			Weight:      requestData.Weight,
			UserID:      user.ID,
		}

		if _, err := s.store.Area().Create(&area); err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, nil, http.StatusCreated)
	}
}
