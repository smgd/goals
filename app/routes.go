package app

func (s *server) routes() {
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/hello", s.privateRoute(s.hello())).Methods("POST")
}
