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
		tokenSigningKey: []byte("blabla"),
		db:              &db,
	}
	s.routes()
	return s
}

func Run() error {
	fmt.Println("listening...")
	db, err := gorm.Open("postgres", "postgres://goals:qwerty@localhost:5432/goals?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	s := newServer(*db)
	return http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, s.router))
}
