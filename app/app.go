package app

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type server struct {
	router          *mux.Router
	tokenSigningKey []byte
}

func newServer() *server {
	s := &server{
		router:          mux.NewRouter(),
		tokenSigningKey: []byte("blabla"),
	}
	s.routes()
	return s
}

func Run() error {
	fmt.Println("listening...")
	s := newServer()
	return http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, s.router))
}
