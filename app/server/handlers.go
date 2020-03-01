package server

import (
	"net/http"

	"github.com/smgd/goals/app/models"

	"github.com/jinzhu/copier"
)

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

		if err := s.decodeBody(r, &requestData); err != nil {
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

		if err := s.decodeBody(r, &requestData); err != nil {
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
		user := s.getRequestUser(r)
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
	type goalsResponse struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type areasResponse struct {
		ID          uint            `json:"id"`
		Weight      int             `json:"weight"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Icon        string          `json:"icon"`
		IsFavourite bool            `json:"is_favourite"`
		Goals       []goalsResponse `json:"goals"`
	}
	type response struct {
		Areas []areasResponse `json:"areas"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.getRequestUser(r)

		userAreas, err := s.store.Area().FindAreasByUserID(user.ID)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
		}

		resp := response{Areas: []areasResponse{}}

		if err := copier.Copy(&resp.Areas, &userAreas); err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		for i, area := range resp.Areas {
			resp.Areas[i].Goals = []goalsResponse{}

			goals, err := s.store.Goal().FindGoalsByAreaID(area.ID)
			if err != nil {
				s.respond(w, nil, http.StatusBadRequest)
			}

			if err := copier.Copy(&resp.Areas[i].Goals, &goals); err != nil {
				s.respond(w, nil, http.StatusInternalServerError)
				return
			}
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

		if err := s.decodeBody(r, &requestData); err != nil {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		user := s.getRequestUser(r)

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

func (s *Server) handleGetGoals() http.HandlerFunc {
	type goalsResponse struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AreaID      uint   `json:"area_id"`
	}
	type response struct {
		Goals []goalsResponse `json:"goals"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.getRequestUser(r)

		userGoals, err := s.store.Goal().FindGoalsByUserID(user.ID)
		if err != nil {
			s.respond(w, nil, http.StatusBadRequest)
		}

		resp := response{Goals: []goalsResponse{}}

		if err := copier.Copy(&resp.Goals, &userGoals); err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, resp, http.StatusOK)
	}
}
