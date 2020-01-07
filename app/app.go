package app

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"os"
)

type server struct {
	router          *mux.Router
	tokenSigningKey []byte
	db              *gorm.DB
}

func newServer(db gorm.DB) *server {
	s := &server{
		router:          mux.NewRouter(),
		tokenSigningKey: []byte(os.Getenv("TOKEN_SIGNING_KEY")),
		db:              &db,
	}
	s.routes()
	return s
}

func Run() error {
	fmt.Println("listening...")
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
	return http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, s.router))
}
