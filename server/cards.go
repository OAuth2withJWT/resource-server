package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/OAuth2withJWT/resource-server/app"
	"github.com/gorilla/mux"
)

func (s *Server) handleGetTotalBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := s.app.CardService.GetTotalBalanceByUserId(userId)
	if err != nil {
		if invalidUserIdErr, ok := err.(*app.InvalidUserIdError); ok {
			log.Print(invalidUserIdErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balance)
}
func (s *Server) handleGetCards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.app.CardService.GetCardsByUserId(userId)
	if err != nil {
		log.Print("Error retrieving cards: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cards)
}
