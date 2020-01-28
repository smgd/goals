package app

import (
	"fmt"
	"net/http"
)

func (s *server) addRoute(path string, handleFunc func(http.ResponseWriter,
	*http.Request), methods ...string) {
	methods = append(methods, "OPTIONS")
	s.router.HandleFunc(fmt.Sprintf("/api/%s", path), s.withCORS(handleFunc)).Methods(methods...)
}

func (s *server) addPrivateRoute(path string, handleFunc func(http.ResponseWriter,
	*http.Request), methods ...string) {
	s.addRoute(path, s.privateRoute(handleFunc), methods...)
}

func (s *server) routes() {
	s.addRoute("login", s.handleLogin(), "POST")
	s.addRoute("register", s.handleRegister(), "POST")
	s.addRoute("ping", s.handlePing(), "GET")
	s.addPrivateRoute("whoami", s.handleWhoAmI(), "GET")
	s.addPrivateRoute("areas", s.handleGetAreas(), "GET")
	s.addPrivateRoute("areas", s.handleCreateAreas(), "POST")
}
