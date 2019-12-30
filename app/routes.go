package app

func (s *server) routes() {
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
}
