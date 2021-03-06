package controllers

import "github.com/fegroders/vineosAPI/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Wines routes
	s.Router.HandleFunc("/wines", middlewares.SetMiddlewareJSON(s.CreateWine)).Methods("POST")
	s.Router.HandleFunc("/wines", middlewares.SetMiddlewareJSON(s.GetWines)).Methods("GET")
	s.Router.HandleFunc("/wines/{id}", middlewares.SetMiddlewareJSON(s.GetWine)).Methods("GET")
	s.Router.HandleFunc("/wines/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateWine))).Methods("PUT")
	s.Router.HandleFunc("/wines/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteWine)).Methods("DELETE")
}