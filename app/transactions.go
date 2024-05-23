package app

import (
	"time"
)

type TransactionService struct {
	repository TransactionRepository
}

func NewTransactionService(tr TransactionRepository) *TransactionService {
	return &TransactionService{
		repository: tr,
	}
}

type TransactionRepository interface {
	GetTransactionByCategoryAmountAndDate(category string, amount float64, date time.Time) (Transaction, error)
}

type Transaction struct {
	Id       int       `json:"id"`
	CardId   int       `json:"card_id"`
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
}

func (s *TransactionService) GetTransactionByCategoryAmountAndDate(category string, amount float64, date time.Time) (Transaction, error) {
	transaction, err := s.repository.GetTransactionByCategoryAmountAndDate(category, amount, date)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}
