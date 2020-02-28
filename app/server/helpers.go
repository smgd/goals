package server

import (
	"encoding/json"
	"goals/app/models"
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
