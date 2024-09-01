package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/OAuth2withJWT/resource-server/app"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (tr *TransactionRepository) GetCategoryExpensesByTime(cardId int, time time.Time) ([]app.CategoryExpenseResponse, error) {
	var expenses []app.CategoryExpenseResponse
	rows, err := tr.db.Query("SELECT expense_category, COALESCE(SUM(amount), 0) AS total_spent FROM transactions WHERE card_id = $1 AND transaction_type = 'expense' AND time > $2 GROUP BY expense_category", cardId, time)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense app.CategoryExpenseResponse
		if err := rows.Scan(&expense.ExpenseCategory, &expense.TotalSpent); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (tr *TransactionRepository) GetTransactionsByTime(cardId int, time time.Time) ([]app.Transaction, error) {
	rows, err := tr.db.Query("SELECT id, card_id, TO_CHAR(time, 'DD/MM/YYYY hh:mm:ss'), amount, expense_category, transaction_type, location, destination_account_id, source_account_id FROM transactions WHERE card_id = $1 AND time > $2", cardId, time)
	if err != nil {
		return []app.Transaction{}, err
	}

	var transactions []app.Transaction
	for rows.Next() {
		var transaction app.Transaction
		err := rows.Scan(&transaction.Id, &transaction.CardId, &transaction.Time, &transaction.Amount, &transaction.ExpenseCategory, &transaction.TransactionType, &transaction.Location, &transaction.DestinationAccountId, &transaction.SourceAccountId)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (tr *TransactionRepository) GetTotalAmountByCategoryAndTime(cardId int, category string, time time.Time) (app.AmountResponse, error) {
	var amount app.AmountResponse
	var totalAmount sql.NullFloat64

	err := tr.db.QueryRow("SELECT SUM(amount) AS total_amount FROM transactions WHERE card_id = $1 AND expense_category = $2 AND time > $3", cardId, category, time).Scan(&totalAmount)
	if err != nil {
		return app.AmountResponse{}, err
	}

	if totalAmount.Valid {
		amount.TotalValue = totalAmount.Float64
	} else {
		amount.TotalValue = 0
	}

	return amount, nil
}

func (tr *TransactionRepository) GetTotalAmountByTime(cardId int, time time.Time, transactionType string) (app.AmountResponse, error) {
	var amount app.AmountResponse
	var totalAmount sql.NullFloat64

	err := tr.db.QueryRow("SELECT SUM(amount) AS total_amount FROM transactions WHERE card_id = $1 AND time > $2 AND transaction_type = $3", cardId, time, transactionType).Scan(&totalAmount)
	if err != nil {
		return app.AmountResponse{}, err
	}

	if totalAmount.Valid {
		amount.TotalValue = totalAmount.Float64
	} else {
		amount.TotalValue = 0
	}

	return amount, nil
}

func (tr *TransactionRepository) GetTransactionsByCategoryAndTime(cardId int, category string, time time.Time) ([]app.Transaction, error) {
	rows, err := tr.db.Query("SELECT id, card_id, time, amount, category, location FROM transactions WHERE card_id = $1 AND category = $2 AND time > $3", cardId, category, time)
	if err != nil {
		return []app.Transaction{}, err
	}

	var transactions []app.Transaction
	for rows.Next() {
		var transaction app.Transaction
		err := rows.Scan(&transaction.Id, &transaction.CardId, &transaction.Time, &transaction.Amount, &transaction.ExpenseCategory, &transaction.TransactionType, &transaction.Location, &transaction.DestinationAccountId, &transaction.SourceAccountId)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
