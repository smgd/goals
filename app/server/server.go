package server

import (
	. "goals/app/models"
	"goals/app/store"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// API server
type Server struct {
	config          *Config,
	logger          *logrus.Logger
	router          *mux.Router
	store			*store.Store
	tokenSigningKey []byte
}


// Server constructor
func New(config *Config, db gorm.DB) *Server {
	s := &Server{{
		config:          config,
		logger:			 logrus.New(),
		router:          mux.NewRouter(),
		tokenSigningKey: []byte(os.Getenv("TOKEN_SIGNING_KEY")),
		db:              &db,
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

func (s *Server) configureStore() {
	st := store.New(s.config.Store)

	if err := st.Open(); err != nil {
		return err
	}
}

// Start server
func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info('starting server...')

	dbAddr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open("postgres", dbAddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"package":  "app",
			"function": "gorm.Open",
			"error":    err,
			"dbAddr":   dbAddr,
		}).Error("Can't open database connection")
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(
		&User{},
		&Area{},
		&Goal{},
		&Task{},
	)

	s := newServer(*db)
	return http.ListenAndServe(s.config.BindAddr, handlers.LoggingHandler(os.Stdout, s.router))
}