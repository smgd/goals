package server

import (
	"goals/app/store"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// API server
type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// Server constructor
func New(config *Config) *Server {
	s := &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
	return s
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *Server) configureRouter() {
	s.addRoute("login", s.handleLogin(), "POST")
	s.addRoute("register", s.handleRegister(), "POST")
	s.addRoute("ping", s.handlePing(), "GET")
	s.addPrivateRoute("whoami", s.handleWhoAmI(), "GET")
	s.addPrivateRoute("areas", s.handleGetAreas(), "GET")
	s.addPrivateRoute("areas", s.handleCreateAreas(), "POST")
}

func (s *Server) configureStore() error {
	st := store.New(s.config.Store)

	s.store = st

	return nil
}

// Start server
func Start(config *Config) error {
	s := New(config)

	// configure logger
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Logger configured")

	// configure router
	s.configureRouter()
	s.logger.Info("Router configured")

	// configure store
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Store configured")

	// create connection with database
	if err := s.store.Open(); err != nil {
		return err
	}
	s.logger.Info("Store connected")
	defer s.store.Close()

	s.logger.Info("starting server...")

	return http.ListenAndServe(s.config.BindAddr, handlers.LoggingHandler(os.Stdout, s.router))
}
