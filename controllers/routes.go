package controllers

import "github.com/ntorres0612/ionix-crud/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/signup", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")

	//Drugs routes
	s.Router.HandleFunc("/drugs", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateDrug))).Methods("POST")
	s.Router.HandleFunc("/drugs", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetDrugs))).Methods("GET")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetDrug))).Methods("GET")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateDrug))).Methods("PUT")
	s.Router.HandleFunc("/drugs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteDrug)).Methods("DELETE")

	//Vaccination routes
	s.Router.HandleFunc("/vaccination", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateVaccination))).Methods("POST")
	s.Router.HandleFunc("/vaccination", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetVaccinations))).Methods("GET")
	s.Router.HandleFunc("/vaccination/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetVaccination))).Methods("GET")
	s.Router.HandleFunc("/vaccination/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateVaccination))).Methods("PUT")
	s.Router.HandleFunc("/vaccination/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteVaccination)).Methods("DELETE")
}
