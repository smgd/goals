package app

func (s *server) routes() {
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/whoami", s.privateRoute(s.handlerWhoAmI())).Methods("GET")
}
