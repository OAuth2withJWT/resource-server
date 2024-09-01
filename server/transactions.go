package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/OAuth2withJWT/resource-server/app"
	"github.com/gorilla/mux"
)

func (s *Server) handleGetTotalAmount(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	transactionType := r.URL.Query().Get("transaction_type")

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time, err := time.Parse("2006-01-02 15:04:05", r.URL.Query().Get("date")+" "+r.URL.Query().Get("time"))
	if err != nil {
		log.Print("Invalid date/time: ", time)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.app.CardService.GetCardsByUserId(userId)
	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(cards) == 0 {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var amount app.AmountResponse
	if category == "" {
		amount, err = s.app.TransactionService.GetTotalAmountByTime(cards, time, transactionType)
	} else {
		amount, err = s.app.TransactionService.GetTotalAmountByCategoryAndTime(cards, category, time)
	}
	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(amount)
}

func (s *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		log.Print("Invalid date/time: ", time)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.app.CardService.GetCardsByUserId(userId)
	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(cards) == 0 {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var transactions []app.Transaction
	if category == "" {
		transactions, err = s.app.TransactionService.GetTransactionsByTime(cards, time)
	} else {
		transactions, err = s.app.TransactionService.GetTransactionsByCategoryAndTime(cards, category, time)
	}
	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func (s *Server) handleGetCategoryExpenses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time, err := time.Parse("2006-01-02 15:04:05", r.URL.Query().Get("date")+" "+r.URL.Query().Get("time"))
	if err != nil {
		log.Print("Invalid date/time: ", time)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.app.CardService.GetCardsByUserId(userId)
	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(cards) == 0 {
		log.Print("Invalid user id: ", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var result []app.CategoryExpenseResponse
	result, err = s.app.TransactionService.GetCategoryExpensesByTime(cards, time)

	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
