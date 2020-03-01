package server

func (s *Server) configureRouter() {
	s.addRoute("login", s.handleLogin(), "POST")
	s.addRoute("register", s.handleRegister(), "POST")
	s.addRoute("ping", s.handlePing(), "GET")
	s.addPrivateRoute("whoami", s.handleWhoAmI(), "GET")
	s.addPrivateRoute("areas", s.handleGetAreas(), "GET")
	s.addPrivateRoute("areas", s.handleCreateAreas(), "POST")
	s.addPrivateRoute("goals", s.handleGetGoals(), "GET")
}
