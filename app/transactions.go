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
	GetTotalAmountByTime(cardId int, time time.Time, transactionType string) (AmountResponse, error)
	GetCategoryExpensesByTime(cardId int, time time.Time) ([]CategoryExpenseResponse, error)
}

type Transaction struct {
	Id              int     `json:"id"`
	CardId          int     `json:"card_id"`
	Time            string  `json:"time"`
	Amount          float64 `json:"amount"`
	ExpenseCategory string  `json:"expense_category"`
	TransactionType string  `json:"transaction_type"`
	Location        *string `json:"location"`
}

type AmountResponse struct {
	Category   string  `json:"category"`
	TotalValue float64 `json:"total_amount"`
}

type CategoryExpenseResponse struct {
	ExpenseCategory string  `json:"expense_category"`
	TotalSpent      float64 `json:"total_spent"`
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

func (s *TransactionService) GetTotalAmountByTime(cards []Card, date time.Time, transactionType string) (AmountResponse, error) {
	totalAmount := AmountResponse{Category: "none"}
	for _, card := range cards {
		amount, err := s.repository.GetTotalAmountByTime(card.Id, date, transactionType)
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

func (s *TransactionService) GetCategoryExpensesByTime(cards []Card, date time.Time) ([]CategoryExpenseResponse, error) {
	categoryTotals := make(map[string]float64)

	for _, card := range cards {
		result, err := s.repository.GetCategoryExpensesByTime(card.Id, date)
		if err != nil {
			return nil, err
		}
		for _, expense := range result {
			categoryTotals[expense.ExpenseCategory] += expense.TotalSpent
		}
	}

	var responses []CategoryExpenseResponse
	for category, total := range categoryTotals {
		responses = append(responses, CategoryExpenseResponse{
			ExpenseCategory: category,
			TotalSpent:      total,
		})
	}

	return responses, nil
}
