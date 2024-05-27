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

func (tr *TransactionRepository) GetTransactionsByCategoryAndTime(category string, time time.Time) ([]app.Transaction, error) {
	rows, err := tr.db.Query("SELECT id, card_id, time, amount, category, location FROM transactions WHERE category = $1 AND time > $2", category, time)
	if err != nil {
		return []app.Transaction{}, err
	}

	var transactions []app.Transaction
	for rows.Next() {
		var transaction app.Transaction
		err := rows.Scan(&transaction.Id, &transaction.CardId, &transaction.Time, &transaction.Amount, &transaction.Category, &transaction.Location)
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

	err := tr.db.QueryRow("SELECT SUM(amount) AS total_amount FROM transactions WHERE card_id = $1 AND category = $2 AND time > $3", cardId, category, time).Scan(&totalAmount)
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

func (tr *TransactionRepository) GetTotalAmountByTime(cardId int, time time.Time) (app.AmountResponse, error) {
	var amount app.AmountResponse
	var totalAmount sql.NullFloat64

	err := tr.db.QueryRow("SELECT SUM(amount) AS total_amount FROM transactions WHERE card_id = $1 AND time > $2", cardId, time).Scan(&totalAmount)
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
