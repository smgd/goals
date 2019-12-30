package app

func (s *server) routes() {
	s.router.HandleFunc("/api/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/api/register", s.handleRegister()).Methods("POST")
	s.router.HandleFunc("/api/whoami", s.privateRoute(s.handlerWhoAmI())).Methods("GET")
}
