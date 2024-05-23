package postgres

import (
	"database/sql"
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

func (tr *TransactionRepository) GetTransactionByCategoryAmountAndDate(category string, amount float64, date time.Time) (app.Transaction, error) {
	var transaction app.Transaction
	row := tr.db.QueryRow("SELECT id, card_id, date, amount, category FROM transactions WHERE category = $1 AND amount = $2 AND date = $3 LIMIT 1", category, amount, date)

	err := row.Scan(&transaction.Id, &transaction.CardId, &transaction.Date, &transaction.Amount, &transaction.Category)
	if err != nil {
		return app.Transaction{}, err
	}

	return transaction, nil
}
