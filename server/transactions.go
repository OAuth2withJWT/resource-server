package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) handleRetrieveTransactionsByCriteria(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	transaction, err := s.app.TransactionService.GetTransactionByCategoryAmountAndDate(category, amount, date)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}
