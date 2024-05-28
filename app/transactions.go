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
	GetTransactionsByCategoryAndTime(cardId int, category string, time time.Time) ([]Transaction, error)
	GetTransactionsByTime(cardId int, time time.Time) ([]Transaction, error)
	GetTotalAmountByCategoryAndTime(cardId int, category string, time time.Time) (AmountResponse, error)
	GetTotalAmountByTime(cardId int, time time.Time) (AmountResponse, error)
}

type Transaction struct {
	Id       int       `json:"id"`
	CardId   int       `json:"card_id"`
	Time     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Location string    `json:"location"`
}

type AmountResponse struct {
	Category   string  `json:"category"`
	TotalValue float64 `json:"total_amount"`
}

func (s *TransactionService) GetTotalAmountByCategoryAndTime(cards []Card, category string, date time.Time) (AmountResponse, error) {
	totalAmount := AmountResponse{Category: category}
	for _, card := range cards {
		amount, err := s.repository.GetTotalAmountByCategoryAndTime(card.Id, category, date)
		if err != nil {
			return AmountResponse{}, err
		}

		totalAmount.TotalValue += amount.TotalValue
	}

	return totalAmount, nil
}

func (s *TransactionService) GetTotalAmountByTime(cards []Card, date time.Time) (AmountResponse, error) {
	totalAmount := AmountResponse{Category: "none"}
	for _, card := range cards {
		amount, err := s.repository.GetTotalAmountByTime(card.Id, date)
		if err != nil {
			return AmountResponse{}, err
		}

		totalAmount.TotalValue += amount.TotalValue
	}

	return totalAmount, nil
}

func (s *TransactionService) GetTransactionsByCategoryAndTime(cards []Card, category string, date time.Time) ([]Transaction, error) {
	allTransactions := []Transaction{}
	for _, card := range cards {
		transactions, err := s.repository.GetTransactionsByCategoryAndTime(card.Id, category, date)
		if err != nil {
			return []Transaction{}, err
		}

		allTransactions = append(allTransactions, transactions...)
	}

	return allTransactions, nil
}

func (s *TransactionService) GetTransactionsByTime(cards []Card, date time.Time) ([]Transaction, error) {
	allTransactions := []Transaction{}
	for _, card := range cards {
		transactions, err := s.repository.GetTransactionsByTime(card.Id, date)
		if err != nil {
			return []Transaction{}, err
		}

		allTransactions = append(allTransactions, transactions...)
	}

	return allTransactions, nil
}
