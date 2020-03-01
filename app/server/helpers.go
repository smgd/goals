package server

import (
	"encoding/json"
	"fmt"
	"github.com/smgd/goals/app/models"
	"net/http"

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

func (s *Server) getRequestUser(r *http.Request) *models.User {
	return r.Context().Value("User").(*models.User)
}

func (s *Server) decodeBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(&v)
}

func (s *Server) addRoute(path string, handleFunc func(http.ResponseWriter, *http.Request), methods ...string) {
	methods = append(methods, "OPTIONS")
	s.router.HandleFunc(fmt.Sprintf("/api/%s", path), withCORS(handleFunc)).Methods(methods...)
}

func (s *Server) addPrivateRoute(path string, handleFunc func(http.ResponseWriter, *http.Request), methods ...string) {
	s.addRoute(path, s.privateRoute(handleFunc), methods...)
}
