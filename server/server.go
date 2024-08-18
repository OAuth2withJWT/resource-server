package server

import (
	"log"
	"net/http"

	"github.com/OAuth2withJWT/resource-server/app"
	"github.com/OAuth2withJWT/resource-server/config"
	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	app       *app.Application
	RSAConfig config.RSAConfig
}

func New(a *app.Application) *Server {
	s := &Server{
		router:    mux.NewRouter(),
		app:       a,
		RSAConfig: config.LoadRSAConfig(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	log.Println("Server started on port 3000")
	return http.ListenAndServe(":3000", s.router)
}

func (s *Server) setupRoutes() {
	s.router.Handle("/api/cards/balance/{user_id}", s.protected(http.HandlerFunc(s.handleGetTotalBalance))).Methods("GET")
	s.router.Handle("/api/transactions/{user_id}", s.protected(http.HandlerFunc(s.handleGetTransactions))).Methods("GET")
	s.router.Handle("/api/transactions/amount/{user_id}", s.protected(http.HandlerFunc(s.handleGetTotalAmount))).Methods("GET")
}
