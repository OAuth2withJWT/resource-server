package server

import (
	"log"
	"net/http"

	"github.com/OAuth2withJWT/resource-server/app"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	app    *app.Application
}

func New(a *app.Application) *Server {
	s := &Server{
		router: mux.NewRouter(),
		app:    a,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", s.router)
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/api/cards/balance/{user_id}", s.handleGetBalance).Methods("GET")
	s.router.HandleFunc("/api/transactions/search", s.handleRetrieveTransactionsByCriteria).Methods("GET")
}
